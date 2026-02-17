import 'dart:developer';

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:flutter/material.dart';
import 'package:frontend/models/connection_status.dart';
import 'package:frontend/models/profile.dart';
import 'package:frontend/pages/contacts_page.dart';
import 'package:frontend/services/contact_service.dart';
import 'package:provider/provider.dart';

import '../services/auth_service.dart';
import '../services/search_service.dart';

class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  List<Profile> results = [];

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          SizedBox(height: 20),
          FractionallySizedBox(
            widthFactor: 0.8,
            child: TextField(
              onSubmitted: (isbn) async {
                try {

                  final searchService = SearchService(
                    context.read<ContactService?>()!,
                    context.read<AuthService>(),
                    FirebaseFirestore.instance,
                  );
                  final profiles = await searchService.searchByISBN(isbn);
                  setState(() {
                    results = profiles;
                  });
                } catch (e) {
                  log('ERROR: $e');
                }
              },
              decoration: InputDecoration(
                prefixIcon: Icon(Icons.search),
                hintText: 'Enter ISBN and press Enter'
              ),
            ),
          ),
          if (results.isNotEmpty)
            Expanded(
              child: GridView.builder(
                itemCount: results.length,
                gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                  maxCrossAxisExtent: 200,
                  childAspectRatio: 0.7,
                  crossAxisSpacing: 10,
                  mainAxisSpacing: 10,
                ),
                itemBuilder: (context, index) {
                  return ContactCard(
                    contact: results[index],
                    connectionStatus: ConnectionStatus.accepted,
                  );
                },
              ),
            ),
        ],
      ),
    );
  }
}
