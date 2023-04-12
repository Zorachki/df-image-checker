# df-image-checker
Dockerfile images checker helps to find images in Dockerfiles which are not included in your allowed list.

Sometimes your team has unlimited access to public images, and you can't close that availability. But also you'd like them to use only your own images pool. You need for automation, of course.

## How it works
You set the pattern of your images and path to your Dockerfile using flag in CLI or env var. The script checks specified Dockerfile and prints first discovered foreign image.

Example: Your images are stored in Gitlab *registry.gitlab.com/mycompany.de/infrastructure/images/*. This means all your images in Dockerfile have to link to that store. Like 
```
FROM registry.gitlab.com/mycompany.de/infrastructure/images/superimage:5.9-wow
```
So set *regPattern* variable as *registry.gitlab.com/mycompany.de/infrastructure/images/** by "-p" flag

This image check is more helpful in CI jobs.
You can enable HARD CHECK to get error, if the script found not correct image. It helps you to stop job run.
## Install
You can build this checker by your own using *go build*. But also can use ready-made binary file (x86-64 or ARM64) and test this checker with attached Dockerfile.
```
dficheck -p "registry.gitlab.com/mycompany.de/secureimages/*"
```
Environments:

DOCKERFILE_PATH

HARD_CHECK

REG_PATTERN
