import 'package:dio/dio.dart';
import 'package:front_end/data/models/file_info.dart';
import 'package:front_end/data/repository/core/api_client.dart';

class DownloadRepository {
  final ApiClient apiClient;

  DownloadRepository(this.apiClient);

  Future<VideoInfo> extractInfo(String url) async {
    try {
      final response = await apiClient.dio.post('/extract', data: {'url': url});
      return VideoInfo.fromJson(response.data);
    } on DioException catch (e) {
      throw Exception(e.response?.data['error'] ?? e.message);
    }
  }

  Future<String> download(String url, {String quality = 'best'}) async {
    try {
      final response = await apiClient.dio.post(
        '/download',
        data: {'url': url, 'quality': quality},
      );
      return response.data['file_path'] ?? '';
    } on DioException catch (e) {
      throw Exception(e.response?.data['error'] ?? e.message);
    }
  }
}
