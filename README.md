# Shrimpify

Image file size optimization server. Using this for my [PlantFam](https://github.com/samuel-ping/PlantFam-Android) Android app.

Named after my _Echeveria fleur blanc_ succulent, Shrimp.

## Setup

### WSL

- Don't forget to install libvips package as per the instructions in the bimg package README.
  - [Install on Ubuntu/WSL](https://github.com/libvips/libvips/wiki/Build-for-Ubuntu)
- Set `$PORT` environmental variable to port of your choice (default can be 8080).

### Heroku

- Using the Aptfile buildpack to install `libvips` package in Heroku server.
- When pushing an update from the heroku branch (hosted on GitHub) to the main branch (hosted on Heroku), run `git push heroku heroku:main`
