import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { DEBUG } from '@glimmer/env';

/**
 *
 * `ConfigDownloadLink` is a button to download the client config for a Consul
 * cluster. The config is downloaded from the cloud-consul service, which responds
 * with JSON data containing the bytes (base64 strings) of the config files. The
 * component extracts the bytes of a combined config bundle zip archive, builds a
 * file, and then triggers a browser download.
 * It builds the API request using the attributes of the provided Consul cluster.
 *
 * ```
 * <ConfigDownloadButton
 *   @resource={{aConsulCluster}}
 * />
 * ```
 *
 * @class ConfigDownloadButton
 *
 */
export default class ConfigDownloadButtonComponent extends Component {
  @service api;

  /**
   *
   * The HCP resource whose config should be downloaded. For example,
   * a Consul cluster.
   * @argument resource
   * @type {object}
   *
   */

  /**
   *
   * The HCP resource type. For example, "consul".
   * @argument resourceType
   * @type {string}
   *
   */

  get resource() {
    return this.args.resource;
  }

  get resourceType() {
    return this.args.resourceType;
  }

  @action
  async performDownload(clickEvent) {
    clickEvent.preventDefault();

    let zipBase64Data = await this.getZippedConfigData();
    if (DEBUG) {
      //eslint-disable-next-line no-console
      console.log(zipBase64Data);
    }

    let fileName = this.buildFileName();
    let zipBytes = this.base64ToBytes(zipBase64Data);
    let file = new File([zipBytes], fileName, { type: 'application/zip' });
    let fileURL = window.URL.createObjectURL(file);
    this.downloadFile(fileName, fileURL);
  }

  /**
   * Returns the name of the client config file bundle to download.
   */
  buildFileName() {
    return `client_config_bundle_${this.resourceType}_${this.resource.id}.zip`;
  }

  /**
   * Perform a request to get the config file bundle from the control plane.
   * It returns a Zip file data, encoded in base64.
   * This method is responsible for checking the type of the resource and then
   * invoking the appropriate logic. Different HCP resource types will use different
   * API endpoints and possibly different response formats.
   *
   * @returns {string} Returns a base64 encoded string, representing a zip file.
   */
  async getZippedConfigData() {
    switch (this.resourceType.toLowerCase()) {
      case 'consul':
        return this.getZippedConfigDataConsul();
      default:
        throw new Error(`Don't know how to download config for resource type "${this.resourceType}".`);
    }
  }

  /**
   * Perform a request to get the config file bundle from the Consul Service API.
   * It returns a Zip file data, encoded in base64.
   *
   * @returns {string} Returns a base64 encoded string, representing a zip file.
   */
  async getZippedConfigDataConsul() {
    let {
      location: { organizationId, projectId },
      id,
    } = this.resource;

    try {
      let { fileBundleZip } = await this.api.consul.getClientConfig(
        organizationId,
        projectId,
        id,
        null,
        null,
        true
      );
      return fileBundleZip;
    } catch (e) {
      if (DEBUG) {
        //eslint-disable-next-line no-console
        console.error(e);
      }
      throw e;
    }
  }

  /**
   * Converts a base64 encoded string into a byte array.
   *
   * @param {string} b64Data A base64 encoded string.
   * @returns {string} Returns a byte array.
   */
  base64ToBytes(b64Data) {
    let byteCharacters = atob(b64Data);
    let byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    let byteArray = new Uint8Array(byteNumbers);
    return byteArray;
  }

  /**
   * Makes the browser download a named file.
   *
   * @param {string} fileName The name of the downloaded file.
   * @param {string} fileURL The file URL for the file.
   */
  downloadFile(fileName, fileURL) {
    let a = document.createElement('a');
    document.body.appendChild(a);
    a.style = 'display: none';

    a.href = fileURL;
    a.download = fileName;
    a.click();
    window.URL.revokeObjectURL(fileURL);
  }
}
