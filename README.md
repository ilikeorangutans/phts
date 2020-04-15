Goals
* Upload photos
* Automatically generate different renditions
* Select and group photos by photo specific criteria
  * Group photos into albums, photostreams
  * Different groupings with different meta data like a trip, vacation?
* Build themeable frontend
* Access control
* Let users comment, subscribe
* Multiple users with different rights/ACL
* filtering/hiding of exif tags on publication?

TODO
* Have to rething transactions/db; currently we initialize the db objects with a sqlx.DB reference. But that way we cannot enforce from the outside that they run in a DB. :|

## Database Setup

```
CREATE DATABASE phts_dev;
CREATE ROLE phts_dev WITH LOGIN;
GRANT ALL PRIVILEGES ON DATABASE phts_dev TO phts_dev;
```

## Minio setup
needs
- mcli
- mcli host localhost


# Building

## Releases

Needs docker and buildx for cross compiling. We could probably do without but if I ever want to build arm specific code in the container it'll make things easier. To set up a builder:
```
docker run --rm --privileged docker/binfmt:820fdd95a9972a5308930a2bdfb8573dd4447ad3
docker buildx rm arm-builder
docker buildx create --name arm-builder
docker buildx inspect --bootstrap arm-builder
```

To build release binaries and docker images:

```
make dist-all -j4
```


## Stuff
Icons from https://fontawesome.com/
Spinners from https://github.com/tobiasahlin/SpinKit
