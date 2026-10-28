import 'package:equatable/equatable.dart';
import 'package:front_end/data/models/file_info.dart';

enum DownloadStatus { initial, loading, success, failure }

class DownloadState extends Equatable {
  final DownloadStatus status;
  final VideoInfo? videoInfo;
  final String? filePath;
  final String? errorMessage;

  const DownloadState({
    this.status = DownloadStatus.initial,
    this.videoInfo,
    this.filePath,
    this.errorMessage,
  });

  DownloadState copyWith({
    DownloadStatus? status,
    VideoInfo? videoInfo,
    String? filePath,
    String? errorMessage,
  }) {
    return DownloadState(
      status: status ?? this.status,
      videoInfo: videoInfo ?? this.videoInfo,
      filePath: filePath ?? this.filePath,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  @override
  List<Object?> get props => [status, videoInfo, filePath, errorMessage];
}
