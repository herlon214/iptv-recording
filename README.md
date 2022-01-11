# IPTV Recording
Manages ffmpeg process to recording IPTV streams.

You need `ffmpeg` installed: https://www.ffmpeg.org/download.html

```shell
$ iptv-rec local

2022/01/08 04:41:02 Found 1 items to record
2022/01/08 04:41:02 ---------------------
2022/01/08 04:41:02 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:02 -_-> Start recording Live Name
2022/01/08 04:41:03 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:04 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:05 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:06 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:07 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:08 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:09 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:10 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:11 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:12 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:13 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:14 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:15 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:16 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:17 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:18 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:19 Live Name [0 04 * * *] -> [41m20s] > live? true
2022/01/08 04:41:20 ---> Stop recording
```

Example of `recording.yaml`:
```yaml
items:
  - name: Live Name
    url: http://iptv-stream-url/live/user/pwd/7986.ts
    fileName: recording-$date # $date will be transformed to "2022-01-11.1720" (for example)
    folder: /mnt/data
    schedule: "0 04 * * *" # Cron style
    duration: 41m20s
```

Example of systemctl daemon service:

```
[Unit]
Description=Record IPTV streams

[Service]
User=pi
WorkingDirectory=/home/pi # Folder which contains the recording.yaml file
Type=simple
TimeoutStopSec=10
Restart=on-failure
RestartSec=5
ExecStart=/usr/bin/iptv-rec local

[Install]
WantedBy=multi-user.target
```