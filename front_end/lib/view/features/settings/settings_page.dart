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
      appBar: AppBar(title: const Text('Settings'), centerTitle: true),
      body: ListView(
        children: [
          const Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Text(
              'Download',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          ListTile(
            title: const Text('Default Quality'),
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
              'Appearance',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          SwitchListTile(
            title: const Text('Dark Mode'),
            subtitle: const Text('Will only be saved on this page for now'),
            value: _darkMode,
            onChanged: (value) {
              setState(() => _darkMode = value);
            },
          ),
          const Divider(),
          const Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: Text(
              'Server',
              style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey),
            ),
          ),
          ListTile(
            title: const Text('Backend URL'),
            subtitle: Text(_serverUrl),
            trailing: const Icon(Icons.edit),
            onTap: () async {
              final controller = TextEditingController(text: _serverUrl);
              final result = await showDialog<String>(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('Backend URL'),
                  content: TextField(
                    controller: controller,
                    decoration: const InputDecoration(
                      hintText: 'http://localhost:8080',
                      border: OutlineInputBorder(),
                    ),
                    keyboardType: TextInputType.url,
                  ),
                  actions: [
                    TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancel')),
                    FilledButton(
                      onPressed: () => Navigator.pop(ctx, controller.text.trim()),
                      child: const Text('Save'),
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
          const ListTile(title: Text('Social Downloader'), subtitle: Text('Version 1.0.0 — Web')),
          const ListTile(
            title: Text(
              '© nazarbaazaar.com',
              style: TextStyle(fontSize: 12, color: Colors.grey),
              textAlign: TextAlign.center,
            ),
          ),
        ],
      ),
    );
  }
}
