package client


var Commit = func() string {
    // if info, ok := debug.ReadBuildInfo(); ok {
	// 	fmt.Println(info, info.Settings)
    //     for _, setting := range info.Settings {
    //         if setting.Key == "vcs.revision" {
    //             return setting.Value
    //         }
    //     }
    // }
    return ""
}()
