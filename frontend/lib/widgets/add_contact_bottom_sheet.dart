import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class AddContactBottomSheet extends StatelessWidget {
  const AddContactBottomSheet({super.key});

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
            child: SearchBar(),
          )
        ],
      ),

    );
  }

}