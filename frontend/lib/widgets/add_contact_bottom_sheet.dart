import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:frontend/models/connection_status.dart';
import 'package:frontend/pages/contacts_page.dart';
import 'package:frontend/services/contact_service.dart';
import 'package:provider/provider.dart';

import '../models/profile.dart';

/// This widget provides the buttom sheet, offering input fields to add new contacts

class AddContactBottomSheet extends StatefulWidget {
  const AddContactBottomSheet({super.key});

  @override
  State<AddContactBottomSheet> createState() => _AddContactBottomSheetState();
}

class _AddContactBottomSheetState extends State<AddContactBottomSheet> {
  Profile? result;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text('Add a new contact'),
          ),
          SizedBox(height: 20),
          FractionallySizedBox(
            widthFactor: 0.8,
            child: TextField(
              onSubmitted: (uid) async {
                try {
                  final contactService = context.read<ContactService?>();
                  Profile? profile = await contactService
                      ?.fetchProfile(uid)
                      .first;
                  setState(() {
                    result = profile;
                  });
                } catch (e) {
                  log('ERROR; $e');
                }
              },
              decoration: InputDecoration(
                prefixIcon: Icon(Icons.search),
                hintText: 'Enter UID',
              ),
            ),
          ),
          SizedBox(height: 20),
          if (result != null)
            InkWell(
              onTap: () {
                context.read<ContactService?>()?.addContact(result!.uid);
              },
              child: SizedBox(
                width: 200,
                height: 280,
                child: ContactCard(
                  contact: result!,
                  connectionStatus: ConnectionStatus.pending,
                ),
              ),
            ),
        ],
      ),
    );
  }
}
