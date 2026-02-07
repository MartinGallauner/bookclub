import 'dart:convert';

import 'package:frontend/models/book.dart';
import 'package:frontend/services/book_api_service.dart';
import 'package:http/http.dart' as http;

class GoogleBooksAPIService implements BookApiService {
  /// This class provides service methods to call the OpenLibary.org API.

  static const String baseUrl = 'https://www.googleapis.com';

  @override
  Future<Book> searchByISBN(String isbn) async {
    final url = Uri.parse('$baseUrl/books/v1/volumes?q=isbn:$isbn');

    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final volumeInfo = data['items'][0]['volumeInfo'];

        return Book(
            title: volumeInfo['title'],
            authors: List<String>.from(volumeInfo['authors'] ?? []),
            isbn: isbn,
            language: "en",
            coverUrl: volumeInfo['imageLinks']?['thumbnail'] ?? "");
      } else {
        throw Exception('API returned status ${response.statusCode}');
      }
      
    } catch (e) {
      throw Exception('Failed to load book because of $e');
    }
  }

}
