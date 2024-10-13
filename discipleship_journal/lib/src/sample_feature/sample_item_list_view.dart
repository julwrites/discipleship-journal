import 'package:flutter/material.dart';

import '../settings/settings_view.dart';
import '../services/play_integrity_service.dart';
import 'sample_item.dart';
import 'sample_item_details_view.dart';

/// Displays a list of SampleItems.
class SampleItemListView extends StatelessWidget {
  SampleItemListView({
    super.key,
    this.items = const [SampleItem(1), SampleItem(2), SampleItem(3)],
  });

  static const routeName = '/';

  final List<SampleItem> items;
  final PlayIntegrityService _playIntegrityService = PlayIntegrityService();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Sample Items'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              Navigator.restorablePushNamed(context, SettingsView.routeName);
            },
          ),
        ],
      ),
      body: ListView.builder(
        restorationId: 'sampleItemListView',
        itemCount: items.length,
        itemBuilder: (BuildContext context, int index) {
          final item = items[index];

          return ListTile(
            title: Text('SampleItem ${item.id}'),
            leading: const CircleAvatar(
              foregroundImage: AssetImage('assets/images/flutter_logo.png'),
            ),
            onTap: () {
              Navigator.restorablePushNamed(
                context,
                SampleItemDetailsView.routeName,
              );
            }
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          bool isIntegrityVerified = await _playIntegrityService.verifyDeviceIntegrity();
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(isIntegrityVerified ? 'Device integrity verified' : 'Device integrity check failed'),
            ),
          );
        },
        tooltip: 'Verify Device Integrity',
        child: const Icon(Icons.security),
      ),
    );
  }
}
