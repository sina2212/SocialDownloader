import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:front_end/view/features/home/download/download_bloc.dart';
import 'package:front_end/view/features/home/download/download_event.dart';
import 'package:front_end/view/features/home/download/download_state.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final TextEditingController _urlController = TextEditingController();

  @override
  void dispose() {
    _urlController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Social Downloader'), centerTitle: true),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: BlocBuilder<DownloadBloc, DownloadState>(
          builder: (context, state) {
            return Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextField(
                  controller: _urlController,
                  decoration: const InputDecoration(
                    labelText: 'لینک پست',
                    border: OutlineInputBorder(),
                    hintText: 'لینک اینستاگرام، یوتیوب، ایکس و ...',
                  ),
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: ElevatedButton(
                        onPressed: state.status == DownloadStatus.loading
                            ? null
                            : () {
                                context.read<DownloadBloc>().add(
                                  ExtractRequested(_urlController.text.trim()),
                                );
                              },
                        child: const Text('استخراج اطلاعات'),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: state.status == DownloadStatus.loading
                            ? null
                            : () {
                                context.read<DownloadBloc>().add(
                                  DownloadRequested(_urlController.text.trim()),
                                );
                              },
                        child: const Text('دانلود'),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 24),
                if (state.status == DownloadStatus.loading)
                  const Center(child: CircularProgressIndicator()),
                if (state.status == DownloadStatus.failure)
                  Text(
                    state.errorMessage ?? 'خطایی رخ داد',
                    style: const TextStyle(color: Colors.red),
                  ),
                if (state.videoInfo != null) ...[
                  const SizedBox(height: 16),
                  Text(state.videoInfo!.title, style: Theme.of(context).textTheme.titleMedium),
                  const SizedBox(height: 8),
                  Text('آپلودر: ${state.videoInfo!.uploader}'),
                  Text('مدت: ${state.videoInfo!.duration.toStringAsFixed(0)} ثانیه'),
                ],
                if (state.filePath != null) ...[
                  const SizedBox(height: 16),
                  Text(
                    'دانلود شد:\n${state.filePath}',
                    style: const TextStyle(color: Colors.green),
                  ),
                ],
              ],
            );
          },
        ),
      ),
    );
  }
}
