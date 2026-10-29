import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:front_end/data/repository/core/api_client.dart';
import 'package:front_end/data/repository/download_repository.dart';
import 'package:front_end/view/features/home/download/download_bloc.dart';

import 'router/app_router.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final apiClient = ApiClient();
    final repository = DownloadRepository(apiClient);

    return BlocProvider(
      create: (_) => DownloadBloc(repository),
      child: MaterialApp.router(
        locale: const Locale('fa', 'IR'),
        title: 'Social Downloader',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
          useMaterial3: true,
        ),
        routerConfig: appRouter,
      ),
    );
  }
}
