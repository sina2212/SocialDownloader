import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:front_end/data/repository/download_repository.dart';

import 'download_event.dart';
import 'download_state.dart';

class DownloadBloc extends Bloc<DownloadEvent, DownloadState> {
  final DownloadRepository repository;

  DownloadBloc(this.repository) : super(const DownloadState()) {
    on<ExtractRequested>(_onExtract);
    on<DownloadRequested>(_onDownload);
    on<ResetDownload>(_onReset);
  }

  Future<void> _onExtract(ExtractRequested event, Emitter<DownloadState> emit) async {
    emit(state.copyWith(status: DownloadStatus.loading, errorMessage: null));

    try {
      final info = await repository.extractInfo(event.url);
      emit(state.copyWith(status: DownloadStatus.success, videoInfo: info));
    } catch (e) {
      emit(state.copyWith(status: DownloadStatus.failure, errorMessage: e.toString()));
    }
  }

  Future<void> _onDownload(DownloadRequested event, Emitter<DownloadState> emit) async {
    emit(state.copyWith(status: DownloadStatus.loading, errorMessage: null));

    try {
      final path = await repository.download(event.url, quality: event.quality);
      emit(state.copyWith(status: DownloadStatus.success, filePath: path));
    } catch (e) {
      emit(state.copyWith(status: DownloadStatus.failure, errorMessage: e.toString()));
    }
  }

  void _onReset(ResetDownload event, Emitter<DownloadState> emit) {
    emit(const DownloadState());
  }
}
