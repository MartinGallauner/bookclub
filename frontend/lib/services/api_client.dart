import 'package:dio/dio.dart';
import 'package:bookclub_api/bookclub_api.dart';


class ApiClient {

  static const baseUrl = String.fromEnvironment('API_BASE_URL');
  late final Dio _dio;
  String? _token;

  ApiClient() {
    _dio = Dio(BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: Duration(seconds: 30),
        receiveTimeout: Duration(seconds: 30),
        headers: {'Content-Type': 'application/json'},
    ));

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        if (_token != null) {
          options.headers['Authorization'] = 'Bearer $_token';
        }
        handler.next(options);
      }
    ));
  }

  void setToken(String token) {
    _token = token;
  }

  SystemApi get systemApi => SystemApi(_dio, standardSerializers);
  AuthApi get authApi => AuthApi(_dio, standardSerializers);
}