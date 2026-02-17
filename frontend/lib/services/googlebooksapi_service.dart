import 'dart:convert';

import 'package:frontend/models/book.dart';
import 'package:frontend/services/book_api_service.dart';
import 'package:http/http.dart' as http;

class GoogleBooksAPIService implements BookApiService {
  /// This class provides service methods to call the OpenLibary.org API.
  final String apiKey;
  static const String baseUrl = 'https://www.googleapis.com';

  GoogleBooksAPIService({required this.apiKey});

  @override
  Future<Book> searchByISBN(String input) async {

    String bookcode = input.replaceAll('-', '').replaceAll(' ', '').trim();

    final url = Uri.parse('$baseUrl/books/v1/volumes?q=isbn:$bookcode&key=$apiKey');

    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);

        if (data['items'] == null || data['items'].isEmpty) {
          throw Exception('No book found for ISBN $input');
        }

        final volumeInfo = data['items'][0]['volumeInfo'];

        return Book(
            title: volumeInfo['title'],
            authors: List<String>.from(volumeInfo['authors'] ?? []),
            isbn: input,
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
