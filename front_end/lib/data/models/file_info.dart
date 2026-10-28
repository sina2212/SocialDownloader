class VideoInfo {
  final String id;
  final String title;
  final String uploader;
  final double duration;
  final String thumbnail;
  final String description;
  final String webpageUrl;

  VideoInfo({
    required this.id,
    required this.title,
    required this.uploader,
    required this.duration,
    required this.thumbnail,
    required this.description,
    required this.webpageUrl,
  });

  factory VideoInfo.fromJson(Map<String, dynamic> json) {
    return VideoInfo(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      uploader: json['uploader'] ?? '',
      duration: (json['duration'] as num?)?.toDouble() ?? 0,
      thumbnail: json['thumbnail'] ?? '',
      description: json['description'] ?? '',
      webpageUrl: json['webpage_url'] ?? '',
    );
  }
}
