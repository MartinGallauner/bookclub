import 'package:flutter/material.dart';

import '../models/book.dart';

class BookCard extends StatelessWidget {
  final Book book;
  final int colorIndex;

  const BookCard({super.key, required this.book, required this.colorIndex});

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Cover image
          Expanded(
            flex: 3,
            child: Container(
              width: double.infinity,
              color: Colors.primaries[colorIndex % Colors.primaries.length],
            ),
          ),
          // Book details
          Padding(
            padding: EdgeInsets.all(8.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  book.title,
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                SizedBox(height: 4),
                Text(
                  book.authors[0],
                  style: TextStyle(color: Colors.grey[600], fontSize: 12),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
