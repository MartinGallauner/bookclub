import 'package:flutter/material.dart';

import '../services/auth_service.dart';

class TopBar extends StatelessWidget {
  const TopBar({super.key});

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text('BookClub'),
      backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      actions: [
        IconButton(onPressed: AuthService().signOut, icon: Icon(Icons.logout)),
      ],
    );
  }
}
