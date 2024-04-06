package config

// func TestConfigFactory(t *testing.T) {
// 	type args struct {
// 		provider string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "none",
// 			args: struct{ provider string }{
// 				provider: "none",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			conf, err := ConfigFactory(tt.args.provider)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ConfigFactory() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !tt.wantErr {
// 				require.NotEmpty(t, conf.String())
// 			}
// 		})
// 	}
// }
