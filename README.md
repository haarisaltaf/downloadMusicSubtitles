# downloadMusicSubtitles
Downloads song subtitles in located in subdirectories of where the main go file is placed. Current api is looking to have closed down recently so need to revisit and decide which subtitle api to switch to.

### How it works:
Walks through every "entry-point" (ie folders or files) in the same/ subdirectories. If its not a directory, check if the suffix is .flac or .wav then check the metadata, grabbing just the title using regex then sending to an api to grab the corresponding lyrics. Converts the api response from json to text then cocnverts to subtitle file -- need to revisit once api site isnt down/ have switched to a different one.

Written in golang.
