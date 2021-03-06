# Azpiri

A launching images generator for Retropie. It will scan the provided folder for ROM files, and
for each one, will check the provided backgrounds and foregrounds folders for PNG images
with the same name as the ROM, composing an image with them to be used as a launching image for that game.
It's mainly designed to work alongside Lars Muldjord's [SkyScraper](https://github.com/muldjord/skyscraper), using its folder structure by default.

![Example image](example.png)

## Usage

Just run `azpiri` to start generating custom launching images, that will be stored in `<roms folder>/images`.

### Optional parameters

* `-r` or `--roms`: Directory to get background images from. By default it's the same directory the application is running from.
* `-b` or `--backgrounds`: Directory to get background images from. By default it's `<roms folder>/media/screenshots/`.
* `-f` or `--foregrounds`: Directory to get foreground images from. By default it's `<roms folder>/media/marquees/`.

## Options

Certain image generation parameters can be changed by modifying `azpiri.json`:

* `BackgroundBrightness`: Percentage to increase/decrease background image brightness.
* `BackgroundBlur`: Amount of blur to apply to background image.
* `ForegroundWidth`: Width to resize the foreground image in pixels.
* `TargetWidth`: Created image width.
* `TargetHeight`: Created image height.