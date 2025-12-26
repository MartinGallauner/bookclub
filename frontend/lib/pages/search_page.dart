import 'package:flutter/material.dart';

class SearchPage extends StatelessWidget {
  const SearchPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          SizedBox(height: 20),
          SearchBar(
            hintText: 'What book are you looking for?',
          )

        ],

      ),
    );
  }


}