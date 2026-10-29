import 'package:flutter/material.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  String _defaultQuality = 'best';
  bool _darkMode = false;
  String _serverUrl = 'http://localhost:8080';

  final List<String> _qualities = ['best', '1080p', '720p', '480p', '360p', 'worst'];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('تنظیمات'), centerTitle: true),
      body: ListView(
        children: [
          const Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Text(
              'دانلود',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          ListTile(
            title: const Text('کیفیت پیش‌فرض'),
            subtitle: Text(_defaultQuality),
            trailing: DropdownButton<String>(
              value: _defaultQuality,
              underline: const SizedBox.shrink(),
              items: _qualities.map((q) => DropdownMenuItem(value: q, child: Text(q))).toList(),
              onChanged: (value) {
                if (value != null) {
                  setState(() => _defaultQuality = value);
                }
              },
            ),
          ),
          const Divider(),
          const Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Text(
              'ظاهر',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          SwitchListTile(
            title: const Text('حالت تاریک'),
            subtitle: const Text('فعلاً فقط در این صفحه ذخیره می‌شود'),
            value: _darkMode,
            onChanged: (value) {
              setState(() => _darkMode = value);
            },
          ),
          const Divider(),
          const Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Text(
              'سرور',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          ListTile(
            title: const Text('آدرس بک‌اند'),
            subtitle: Text(_serverUrl),
            trailing: const Icon(Icons.edit),
            onTap: () async {
              final controller = TextEditingController(text: _serverUrl);
              final result = await showDialog<String>(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('آدرس سرور'),
                  content: TextField(
                    controller: controller,
                    decoration: const InputDecoration(
                      hintText: 'http://localhost:8080',
                      border: OutlineInputBorder(),
                    ),
                    keyboardType: TextInputType.url,
                  ),
                  actions: [
                    TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('انصراف')),
                    FilledButton(
                      onPressed: () => Navigator.pop(ctx, controller.text.trim()),
                      child: const Text('ذخیره'),
                    ),
                  ],
                ),
              );
              if (result != null && result.isNotEmpty) {
                setState(() => _serverUrl = result);
              }
            },
          ),
          const Divider(),
          const Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Text(
              'درباره',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          const ListTile(title: Text('Social Downloader'), subtitle: Text('نسخه ۱.۰.۰ — وب')),
        ],
      ),
    );
  }
}
