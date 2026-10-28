import 'package:equatable/equatable.dart';

abstract class DownloadEvent extends Equatable {
  const DownloadEvent();

  @override
  List<Object?> get props => [];
}

class ExtractRequested extends DownloadEvent {
  final String url;

  const ExtractRequested(this.url);

  @override
  List<Object?> get props => [url];
}

class DownloadRequested extends DownloadEvent {
  final String url;
  final String quality;

  const DownloadRequested(this.url, {this.quality = 'best'});

  @override
  List<Object?> get props => [url, quality];
}

class ResetDownload extends DownloadEvent {}
