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

## Database Setup

```
CREATE DATABASE phts_dev;
CREATE ROLE phts_dev WITH LOGIN;
GRANT ALL PRIVILEGES ON DATABASE phts_dev TO phts_dev;
```
