import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:front_end/view/features/home/download/download_bloc.dart';
import 'package:front_end/view/features/home/download/download_state.dart';

class HistoryPage extends StatelessWidget {
  const HistoryPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('تاریخچه دانلودها'), centerTitle: true),
      body: BlocBuilder<DownloadBloc, DownloadState>(
        builder: (context, state) {
          final hasHistory = state.filePath != null || state.videoInfo != null;

          if (!hasHistory) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.history, size: 72, color: Theme.of(context).colorScheme.outline),
                  const SizedBox(height: 16),
                  Text(
                    'هنوز دانلودی ثبت نشده است',
                    style: Theme.of(context).textTheme.titleMedium
                        ?.copyWith(color: Theme.of(context).colorScheme.outline),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'پس از دانلود موفق، اینجا نمایش داده می‌شود',
                    style: Theme.of(context).textTheme.bodyMedium
                        ?.copyWith(color: Theme.of(context).colorScheme.outline),
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            );
          }

          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              if (state.videoInfo != null)
                Card(
                  child: ListTile(
                    leading: state.videoInfo!.thumbnail.isNotEmpty
                        ? ClipRRect(
                            borderRadius: BorderRadius.circular(8),
                            child: Image.network(
                              state.videoInfo!.thumbnail,
                              width: 56,
                              height: 56,
                              fit: BoxFit.cover,
                              errorBuilder: (c, o, s) => const Icon(Icons.video_library, size: 40),
                            ),
                          )
                        : const Icon(Icons.video_library, size: 40),
                    title: Text(
                      state.videoInfo!.title,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    subtitle: Text(
                      'آپلودر: ${state.videoInfo!.uploader}\n'
                      'مدت: ${state.videoInfo!.duration.toStringAsFixed(0)} ثانیه',
                    ),
                    isThreeLine: true,
                  ),
                ),
              if (state.filePath != null) ...[
                const SizedBox(height: 12),
                Card(
                  color: Colors.green.shade50,
                  child: ListTile(
                    leading: const Icon(Icons.check_circle, color: Colors.green),
                    title: const Text('آخرین دانلود موفق'),
                    subtitle: Text(state.filePath!),
                  ),
                ),
              ],
            ],
          );
        },
      ),
    );
  }
}
