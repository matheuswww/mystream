package ffmpeg

var resolutions = []struct {
	width, height  int
	outputFile     string
	audioBitrate   string
	videoBitrate   string
}{
	{1920, 1080, "video_1080p.m3u8", "5000k", "5000k"},
	{1280, 720, "video_720p.m3u8", "2500k", "2500k"},
	{854, 480, "video_480p.m3u8", "1200k", "1200k"},
	{640, 360, "video_360p.m3u8", "800k", "800k"},
}

var numberResolutions = len(resolutions)