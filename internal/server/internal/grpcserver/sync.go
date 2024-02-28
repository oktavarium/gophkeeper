package grpcserver

// func (s *GrpcServer) Sync(stream pbapi.GophKeeper_SyncServer) error {
// 	for {
// 		_, err := stream.Recv()
// 		if err == io.EOF {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}
//
// 		if err := stream.Send(&pbapi.SyncResponse{}); err != nil {
// 			return err
// 		}
// 	}
// }
