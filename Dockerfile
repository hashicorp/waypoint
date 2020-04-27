FROM scratch

COPY ./waypoint /waypoint
COPY ./waypoint-entrypoint /waypoint-entrypoint

ENTRYPOINT ["/waypoint"]
