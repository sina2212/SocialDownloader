class VideoInfo {
  final String id;
  final String title;
  final String uploader;
  final double duration;
  final String thumbnail;
  final String description;
  final String webpageUrl;
  final String mediaType;
  final List<FormatInfo> formats;

  VideoInfo({
    required this.id,
    required this.title,
    required this.uploader,
    required this.duration,
    required this.thumbnail,
    required this.description,
    required this.webpageUrl,
    this.mediaType = '',
    this.formats = const [],
  });

  factory VideoInfo.fromJson(Map<String, dynamic> json) {
    final formatsJson = json['formats'] as List<dynamic>? ?? [];
    return VideoInfo(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      uploader: json['uploader'] ?? '',
      duration: (json['duration'] as num?)?.toDouble() ?? 0,
      thumbnail: json['thumbnail'] ?? '',
      description: json['description'] ?? '',
      webpageUrl: json['webpage_url'] ?? '',
      mediaType: json['media_type'] ?? '',
      formats: formatsJson
          .map((e) => FormatInfo.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}

class FormatInfo {
  final String formatId;
  final String ext;
  final String resolution;
  final int filesize;
  final String note;
  final String vcodec;
  final String acodec;
  final double fps;
  final String url;

  FormatInfo({
    required this.formatId,
    required this.ext,
    required this.resolution,
    required this.filesize,
    required this.note,
    this.vcodec = '',
    this.acodec = '',
    this.fps = 0,
    this.url = '',
  });

  factory FormatInfo.fromJson(Map<String, dynamic> json) {
    return FormatInfo(
      formatId: json['format_id'] ?? '',
      ext: json['ext'] ?? '',
      resolution: json['resolution'] ?? '',
      filesize: (json['filesize'] as num?)?.toInt() ?? 0,
      note: json['note'] ?? '',
      vcodec: json['vcodec'] ?? '',
      acodec: json['acodec'] ?? '',
      fps: (json['fps'] as num?)?.toDouble() ?? 0,
      url: json['url'] ?? '',
    );
  }
}
