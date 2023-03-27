#! /usr/bin/bash
function refresh_docker_images {
   # Arguments:
   #   $1 - Path to top level pkg source
   #   $2 - Which make target to invoke (optional)
   #
   # Return:
   #   0 - success
   #   * - failure

   if ! test -d "$1"
   then
      err "ERROR: '$1' is not a directory. refresh_docker_images must be called with the path to the top level source as the first argument'"
      return 1
   fi

   local sdir="$1"
   local targets="$2"

   test -n "${targets}" || targets="docker-images"

   make -C "${sdir}" ${targets}
   return $?
}

function build_ui {
   # Arguments:
   #   $1 - Path to the top level pkg source
   #   $2 - The docker image to run the build within (optional)
   #   $3 - Version override
   #
   # Returns:
   #   0 - success
   #   * - error
   #
   # Notes:
   #   Use the GIT_COMMIT environment variable to pass off to the build

   if ! test -d "$1"
   then
      err "ERROR: '$1' is not a directory. build_ui must be called with the path to the top level source as the first argument'"
      return 1
   fi

   local image_name=${UI_BUILD_CONTAINER_DEFAULT}
   if test -n "$2"
   then
      image_name="$2"
   fi

   local sdir="$1"
   local ui_dir="${1}/ui"

   # parse the version
   version=$(parse_version "${sdir}")

   if test -n "$3"
   then
      version="$3"
   fi

   local commit_hash="${GIT_COMMIT}"
   if test -z "${commit_hash}"
   then
      commit_hash=$(git rev-parse --short HEAD)
   fi
   local logo_type="${PKG_BINARY_TYPE}"
   if test "$logo_type" != "oss"
   then
     logo_type="enterprise"
   fi

   # make sure we run within the ui dir
   pushd ${ui_dir} > /dev/null

   status "Creating the UI Build Container with image: ${image_name} and version '${version}'"
   local container_id=$(docker create -it -e "PKG_GIT_SHA=${commit_hash}" -e "PKG_VERSION=${version}" -e "PKG_BINARY_TYPE=${PKG_BINARY_TYPE}" ${image_name})
   local ret=$?
   if test $ret -eq 0
   then
      status "Copying the source from '${ui_dir}' to /pkg-src within the container"
      (
         tar -c $(ls -A | grep -v "^(node_modules\|dist\|tmp)") | docker cp - ${container_id}:/pkg-src &&
         status "Running build in container" && docker start -i ${container_id} &&
         rm -rf ${1}/ui/dist &&
         status "Copying back artifacts" && docker cp ${container_id}:/pkg-src/dist ${1}/ui/dist
      )
      ret=$?
      docker rm ${container_id} > /dev/null
   fi

   # Check the version is baked in correctly
   if test ${ret} -eq 0
   then
      local ui_vers=$(ui_version "${1}/ui/dist/index.html")
      if test "${version}" != "${ui_vers}"
      then
         err "ERROR: UI version mismatch. Expecting: '${version}' found '${ui_vers}'"
         ret=1
      fi
   fi

   # Check the logo is baked in correctly
   if test ${ret} -eq 0
   then
     local ui_logo_type=$(ui_logo_type "${1}/ui/dist/index.html")
     if test "${logo_type}" != "${ui_logo_type}"
     then
       err "ERROR: UI logo type mismatch. Expecting: '${logo_type}' found '${ui_logo_type}'"
       ret=1
     fi
   fi

   # Copy UI over ready to be packaged into the binary
   if test ${ret} -eq 0
   then
      rm -rf ${1}/ui/dist
      mkdir -p ${1}/pkg
      cp -r ${1}/ui/dist ${1}/ui/dist
   fi

   popd > /dev/null
   return $ret
}

function build_assetfs {
   # Arguments:
   #   $1 - Path to the top level pkg source
   #   $2 - The docker image to run the build within (optional)
   #
   # Returns:
   #   0 - success
   #   * - error
   #
   # Note:
   #   The GIT_COMMIT, GIT_DIRTY and GIT_DESCRIBE environment variables will be used if present

   if ! test -d "$1"
   then
      err "ERROR: '$1' is not a directory. build_assetfs must be called with the path to the top level source as the first argument'"
      return 1
   fi

   local sdir="$1"
   local image_name=${GO_BUILD_CONTAINER_DEFAULT}
   if test -n "$2"
   then
      image_name="$2"
   fi

   pushd ${sdir} > /dev/null
   status "Creating the Go Build Container with image: ${image_name}"
   local container_id=$(docker create -it -e GIT_COMMIT=${GIT_COMMIT} -e GIT_DIRTY=${GIT_DIRTY} -e GIT_DESCRIBE=${GIT_DESCRIBE} ${image_name} make static-assets ASSETFS_PATH=bindata_assetfs.go)
   local ret=$?
   if test $ret -eq 0
   then
      status "Copying the sources from '${sdir}/(ui/dist|GNUmakefile)' to /${PKG_NAME}/pkg"
      (
         tar -c ui/dist GNUmakefile | docker cp - ${container_id}:/${PKG_NAME} &&
         status "Running build in container" && docker start -i ${container_id} &&
         status "Copying back artifacts" && docker cp ${container_id}:/${PKG_NAME}/bindata_assetfs.go ${sdir}/pkg/server/gen/bindata_ui.go
      )
      ret=$?
      docker rm ${container_id} > /dev/null
   fi
   popd >/dev/null
   return $ret
}

function build_pkg_post {
   # Arguments
   #   $1 - Path to the top level pkg source
   #   $2 - Subdirectory under pkg/bin (Optional)
   #
   # Returns:
   #   0 - success
   #   * - error
   #
   # Notes:
   #   pkg/bin is where to place binary packages
   #   pkg.bin.new is where the just built binaries are located
   #   bin is where to place the local systems versions

   if ! test -d "$1"
   then
      err "ERROR: '$1' is not a directory. build_pkg_post must be called with the path to the top level source as the first argument'"
      return 1
   fi

   local sdir="$1"

   local extra_dir_name="$2"
   local extra_dir=""

   if test -n "${extra_dir_name}"
   then
      extra_dir="${extra_dir_name}/"
   fi

   pushd "${sdir}" > /dev/null
   # recreate the pkg dir
   if [ -d "pkg/bin/${extra_dir}" ]; then
      rm -r pkg/bin/${extra_dir}*
   fi
   mkdir -p pkg/bin/${extra_dir} 2> /dev/null

   # move all files in pkg.new into pkg
   cp -r pkg.bin.new/${extra_dir}* pkg/bin/${extra_dir}
   rm -r pkg.bin.new

   DEV_PLATFORM="./pkg/bin/${extra_dir}$(go env GOOS)_$(go env GOARCH)"
   for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f 2>/dev/null)
   do
      # recreate the bin dir
      if [ -d "bin/" ]; then
         rm -r bin/*
      fi
      mkdir -p bin 2> /dev/null

      cp ${F} bin/
      cp ${F} ${MAIN_GOPATH}/bin
   done

   popd > /dev/null

   return 0
}

function build_pkg {
   # Arguments:
   #   $1 - Path to the top level pkg source
   #   $2 - Subdirectory to put binaries in under pkg/bin (optional - must specify if needing to specify the docker image)
   #   $3 - The docker image to run the build within (optional)
   #
   # Returns:
   #   0 - success
   #   * - error
   #
   # Note:
   #   The GOLDFLAGS and GOTAGS environment variables will be used if set
   #   If the PKG_DEV_MODE environment var is truthy only the local platform/architecture is built.
   #   If the XC_OS or the XC_ARCH environment vars are present then only those platforms/architectures
   #   will be built. Otherwise all supported platform/architectures are built

   if ! test -d "$1"
   then
      err "ERROR: '$1' is not a directory. build_pkg must be called with the path to the top level source as the first argument'"
      return 1
   fi

   local sdir="$1"
   local extra_dir_name="$2"
   local extra_dir=""
   local image_name=${GO_BUILD_CONTAINER_DEFAULT}
   if test -n "$3"
   then
      image_name="$3"
   fi

   pushd ${sdir} > /dev/null
   if is_set "${PKG_DEV_MODE}"
   then
      if test -z "${XC_OS}"
      then
         XC_OS=$(go env GOOS)
      fi

      if test -z "${XC_ARCH}"
      then
         XC_ARCH=$(go env GOARCH)
      fi
   fi
   XC_OS=${XC_OS:-"solaris darwin freebsd linux windows"}
   XC_ARCH=${XC_ARCH:-"386 amd64 arm arm64"}

   if test -n "${extra_dir_name}"
   then
      extra_dir="${extra_dir_name}/"
   fi

   # figure out if the compiler supports modules
   local use_modules=0
   if go help modules >/dev/null 2>&1
   then
      use_modules=1
   elif test -n "${GO111MODULE}"
   then
      use_modules=1
   fi

   local volume_mount=
   if is_set "${use_modules}"
   then
      status "Ensuring Go modules are up to date"
      # ensure our go module cache is correct
      go_mod_assert || return 1
      # setup to bind mount our hosts module cache into the container
      volume_mount="--mount=type=bind,source=${MAIN_GOPATH}/pkg/mod,target=/go/pkg/mod"
   fi

   status "Creating the Go Build Container with image: ${image_name}"
   local container_id=$(docker create -it \
      ${volume_mount} \
      -e GOLDFLAGS="${GOLDFLAGS}" \
      -e GOTAGS="${GOTAGS}" \
      ${image_name} \
      ./build-support/scripts/build-local.sh -o "${XC_OS}" -a "${XC_ARCH}")
   ret=$?

   if test $ret -eq 0
   then
      status "Copying the source from '${sdir}' to /${PKG_NAME}"
      (
         tar -c $(ls | grep -v "^(ui\|website\|bin\|pkg\|.git)") | docker cp - ${container_id}:/${PKG_NAME} &&
         status "Running build in container" &&
         docker start -i ${container_id} &&
         status "Copying back artifacts" &&
         docker cp ${container_id}:/${PKG_NAME}/pkg/bin pkg.bin.new
      )
      ret=$?
      docker rm ${container_id} > /dev/null

      if test $ret -eq 0
      then
         build_pkg_post "${sdir}" "${extra_dir_name}"
         ret=$?
      else
         rm -r pkg.bin.new 2> /dev/null
      fi
   fi
   popd > /dev/null
   return $ret
}

function build_pkg_local {
   # Arguments:
   #   $1 - Path to the top level pkg source
   #   $2 - Space separated string of OSes to build. If empty will use env vars for determination.
   #   $3 - Space separated string of architectures to build. If empty will use env vars for determination.
   #   $4 - Subdirectory to put binaries in under pkg/bin (optional)
   #
   # Returns:
   #   0 - success
   #   * - error
   #
   # Note:
   #   The GOLDFLAGS and GOTAGS environment variables will be used if set
   #   If the PKG_DEV_MODE environment var is truthy only the local platform/architecture is built.
   #   If the XC_OS or the XC_ARCH environment vars are present then only those platforms/architectures
   #   will be built. Otherwise all supported platform/architectures are built
   #   The GOXPARALLEL environment variable is used if set

   if ! test -d "$1"
   then
      err "ERROR: '$1' is not a directory. build_pkg must be called with the path to the top level source as the first argument'"
      return 1
   fi

   local sdir="$1"
   local build_os="$2"
   local build_arch="$3"
   local extra_dir_name="$4"
   local extra_dir=""

   if test -n "${extra_dir_name}"
   then
      extra_dir="${extra_dir_name}/"
   fi

   pushd ${sdir} > /dev/null
   if is_set "${PKG_DEV_MODE}"
   then
      if test -z "${XC_OS}"
      then
         XC_OS=$(go env GOOS)
      fi

      if test -z "${XC_ARCH}"
      then
         XC_ARCH=$(go env GOARCH)
      fi
   fi
   XC_OS=${XC_OS:-"solaris darwin freebsd linux windows"}
   XC_ARCH=${XC_ARCH:-"386 amd64 arm arm64"}

   if test -z "${build_os}"
   then
      build_os="${XC_OS}"
   fi

   if test -z "${build_arch}"
   then
      build_arch="${XC_ARCH}"
   fi

   status_stage "==> Building pkg - OSes: ${build_os}, Architectures: ${build_arch}"
   mkdir pkg.bin.new 2> /dev/null

   status "Building sequentially with go build and go-bindata for assets"
   for os in ${build_os}
   do
      for arch in ${build_arch}
      do
         outdir="pkg.bin.new/${extra_dir}${os}_${arch}"
         osarch="${os}/${arch}"

         case "${os}" in
            "darwin" )
               # Only support arm64 for arm darwin
               if test "${arch}" == "arm"
               then
                  continue
               fi
               ;;
            "windows" )
               # Do not build ARM binaries for Windows
               if test "${arch}" == "arm" -o "${arch}" == "arm64"
               then
                  continue
               fi
               ;;
            "freebsd" )
               # Do not build ARM binaries for FreeBSD
               if test "${arch}" == "arm" -o "${arch}" == "arm64"
               then
                  continue
               fi
               ;;
            "solaris" )
               # Only build amd64 for Solaris
               if test "${arch}" != "amd64"
               then
                  continue
               fi
               ;;
            "linux" )
               # build all the binaries for Linux
               ;;
            *)
               continue
            ;;
         esac

         echo "--->   ${osarch}"

         mkdir -p "${outdir}"
         GOBIN_EXTRA=""
         if test "${os}" != "$(go env GOHOSTOS)" -o "${arch}" != "$(go env GOHOSTARCH)"
         then
            GOBIN_EXTRA="${os}_${arch}/"
         fi

         binname="${PKG_NAME}"
         entrypoint="${PKG_SECONDARY_NAME}"
         if [ "$os" == "windows" ];then
               binname="${PKG_NAME}.exe"
         fi

      	GOLDFLAGS+=" -X ${GIT_IMPORT}.Version=${RELEASE_VERSION} -X ${GIT_IMPORT}.VersionPrerelease=${PRERELEASE_VERSION} -X ${GIT_IMPORT}.VersionMetadata=${VERSION_METADATA}"
         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${sdir}/internal/assets/ceb/ceb" "${sdir}/cmd/waypoint-entrypoint"
         CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "${sdir}/internal/assets/ceb/ceb-arm64" "${sdir}/cmd/waypoint-entrypoint"
         # go-bindata requires relative pathing
         pushd "${sdir}/internal/assets" && go-bindata -pkg assets -o prod.go -tags assetsembedded ./ceb && popd
         CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -ldflags "${GOLDFLAGS}" -tags "${GOTAGS}" -o "${outdir}/${binname}" "${sdir}/cmd/waypoint"
         if [ "$os" != "windows" ];then
            CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -tags "${GOTAGS}" -o "${outdir}/${entrypoint}" "${sdir}/cmd/waypoint-entrypoint"
         fi
         #debug_run env GOOS=${os} GOARCH=${arch} go install -ldflags "${GOLDFLAGS}" -tags "${GOTAGS}" && cp "${MAIN_GOPATH}/bin/${GOBIN_EXTRA}${binname}" "${outdir}/${binname}"
         if test $? -ne 0
         then
            err "ERROR: Failed to build pkg for ${osarch}"
            rm -r pkg.bin.new
            return 1
         fi
      done
   done

   build_pkg_post "${sdir}" "${extra_dir_name}"
   if test $? -ne 0
   then
      err "ERROR: Failed postprocessing pkg binaries"
      return 1
   fi
   return 0
}
