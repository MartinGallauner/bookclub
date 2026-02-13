import 'connection_status.dart';

class Connection {
  final String? uid;
  final String requestedBy;
  final DateTime requestedAt;
  final DateTime? acceptedAt;
  final DateTime? rejectedAt;
  final List<String> users;
  final ConnectionStatus status;

  const Connection({
    this.uid,
    required this.requestedBy,
    required this.requestedAt,
    this.acceptedAt,
    this.rejectedAt,
    required this.users,
    required this.status,
  });

  Map<String, dynamic> toMap() {
    return {
      'uid': uid,
      'requestedBy': requestedBy,
      'requestedAt': requestedAt,
      'rejectedAt': rejectedAt,
      'acceptedAt': acceptedAt,
      'users': users,
      'status': status.name,
    };
  }
}
