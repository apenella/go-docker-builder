package response

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {

	var w bytes.Buffer

	input := `{"status":"The push refers to repository [registry.go-docker-builder.test/dummy-image-layers]"}
	{"status":"Preparing","progressDetail":{},"id":"38be7762a5d3"}
	{"status":"Preparing","progressDetail":{},"id":"6a996c0ce279"}
	{"status":"Preparing","progressDetail":{},"id":"d6f45f2d1604"}
	{"status":"Preparing","progressDetail":{},"id":"8407c4f3604d"}
	{"status":"Preparing","progressDetail":{},"id":"4367a98dd925"}
	{"status":"Preparing","progressDetail":{},"id":"36b45d63da70"}
	{"status":"Pushing","progressDetail":{"current":360960,"total":33554431},"progress":"[\u003e                                                  ]    361kB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":360960,"total":33554431},"progress":"[\u003e                                                  ]    361kB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":360960,"total":33554431},"progress":"[\u003e                                                  ]    361kB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":360960,"total":33554431},"progress":"[\u003e                                                  ]    361kB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":360960,"total":33554431},"progress":"[\u003e                                                  ]    361kB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":2163200,"total":33554431},"progress":"[===\u003e                                               ]  2.163MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":2163200,"total":33554431},"progress":"[===\u003e                                               ]  2.163MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":2163200,"total":33554431},"progress":"[===\u003e                                               ]  2.163MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":2163200,"total":33554431},"progress":"[===\u003e                                               ]  2.163MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":2163200,"total":33554431},"progress":"[===\u003e                                               ]  2.163MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":4325888,"total":33554431},"progress":"[======\u003e                                            ]  4.326MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":3965440,"total":33554431},"progress":"[=====\u003e                                             ]  3.965MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":4325888,"total":33554431},"progress":"[======\u003e                                            ]  4.326MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":4325888,"total":33554431},"progress":"[======\u003e                                            ]  4.326MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":4325888,"total":33554431},"progress":"[======\u003e                                            ]  4.326MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":6128128,"total":33554431},"progress":"[=========\u003e                                         ]  6.128MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":6488576,"total":33554431},"progress":"[=========\u003e                                         ]  6.489MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":6128128,"total":33554431},"progress":"[=========\u003e                                         ]  6.128MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":6128128,"total":33554431},"progress":"[=========\u003e                                         ]  6.128MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":6488576,"total":33554431},"progress":"[=========\u003e                                         ]  6.489MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":7930368,"total":33554431},"progress":"[===========\u003e                                       ]   7.93MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":7930368,"total":33554431},"progress":"[===========\u003e                                       ]   7.93MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":8290816,"total":33554431},"progress":"[============\u003e                                      ]  8.291MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":7930368,"total":33554431},"progress":"[===========\u003e                                       ]   7.93MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":8290816,"total":33554431},"progress":"[============\u003e                                      ]  8.291MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":9732608,"total":33554431},"progress":"[==============\u003e                                    ]  9.733MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":10093056,"total":33554431},"progress":"[===============\u003e                                   ]  10.09MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":9732608,"total":33554431},"progress":"[==============\u003e                                    ]  9.733MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":9732608,"total":33554431},"progress":"[==============\u003e                                    ]  9.733MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":10093056,"total":33554431},"progress":"[===============\u003e                                   ]  10.09MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":11534848,"total":33554431},"progress":"[=================\u003e                                 ]  11.53MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":11534848,"total":33554431},"progress":"[=================\u003e                                 ]  11.53MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":12255744,"total":33554431},"progress":"[==================\u003e                                ]  12.26MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":11534848,"total":33554431},"progress":"[=================\u003e                                 ]  11.53MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":12255744,"total":33554431},"progress":"[==================\u003e                                ]  12.26MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":13337088,"total":33554431},"progress":"[===================\u003e                               ]  13.34MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":13337088,"total":33554431},"progress":"[===================\u003e                               ]  13.34MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":14057984,"total":33554431},"progress":"[====================\u003e                              ]  14.06MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":13337088,"total":33554431},"progress":"[===================\u003e                               ]  13.34MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":14057984,"total":33554431},"progress":"[====================\u003e                              ]  14.06MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":14778880,"total":33554431},"progress":"[======================\u003e                            ]  14.78MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":14778880,"total":33554431},"progress":"[======================\u003e                            ]  14.78MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":15860224,"total":33554431},"progress":"[=======================\u003e                           ]  15.86MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":14778880,"total":33554431},"progress":"[======================\u003e                            ]  14.78MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":15499776,"total":33554431},"progress":"[=======================\u003e                           ]   15.5MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":16581120,"total":33554431},"progress":"[========================\u003e                          ]  16.58MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":16581120,"total":33554431},"progress":"[========================\u003e                          ]  16.58MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":17302016,"total":33554431},"progress":"[=========================\u003e                         ]   17.3MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":16581120,"total":33554431},"progress":"[========================\u003e                          ]  16.58MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":17302016,"total":33554431},"progress":"[=========================\u003e                         ]   17.3MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":18383360,"total":33554431},"progress":"[===========================\u003e                       ]  18.38MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":18383360,"total":33554431},"progress":"[===========================\u003e                       ]  18.38MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":19104256,"total":33554431},"progress":"[============================\u003e                      ]   19.1MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":18383360,"total":33554431},"progress":"[===========================\u003e                       ]  18.38MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":19104256,"total":33554431},"progress":"[============================\u003e                      ]   19.1MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":20185600,"total":33554431},"progress":"[==============================\u003e                    ]  20.19MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":20546048,"total":33554431},"progress":"[==============================\u003e                    ]  20.55MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":19825152,"total":33554431},"progress":"[=============================\u003e                     ]  19.83MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":19825152,"total":33554431},"progress":"[=============================\u003e                     ]  19.83MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":20546048,"total":33554431},"progress":"[==============================\u003e                    ]  20.55MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":21627392,"total":33554431},"progress":"[================================\u003e                  ]  21.63MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":21987840,"total":33554431},"progress":"[================================\u003e                  ]  21.99MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":21266944,"total":33554431},"progress":"[===============================\u003e                   ]  21.27MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":21266944,"total":33554431},"progress":"[===============================\u003e                   ]  21.27MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":21987840,"total":33554431},"progress":"[================================\u003e                  ]  21.99MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":23069184,"total":33554431},"progress":"[==================================\u003e                ]  23.07MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":23429632,"total":33554431},"progress":"[==================================\u003e                ]  23.43MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":22708736,"total":33554431},"progress":"[=================================\u003e                 ]  22.71MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":22708736,"total":33554431},"progress":"[=================================\u003e                 ]  22.71MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":23429632,"total":33554431},"progress":"[==================================\u003e                ]  23.43MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":24871424,"total":33554431},"progress":"[=====================================\u003e             ]  24.87MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":24871424,"total":33554431},"progress":"[=====================================\u003e             ]  24.87MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":24510976,"total":33554431},"progress":"[====================================\u003e              ]  24.51MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":24150528,"total":33554431},"progress":"[===================================\u003e               ]  24.15MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":24871424,"total":33554431},"progress":"[=====================================\u003e             ]  24.87MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":26313216,"total":33554431},"progress":"[=======================================\u003e           ]  26.31MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":26313216,"total":33554431},"progress":"[=======================================\u003e           ]  26.31MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":25952768,"total":33554431},"progress":"[======================================\u003e            ]  25.95MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":25592320,"total":33554431},"progress":"[======================================\u003e            ]  25.59MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":26313216,"total":33554431},"progress":"[=======================================\u003e           ]  26.31MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":27755008,"total":33554431},"progress":"[=========================================\u003e         ]  27.76MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":27755008,"total":33554431},"progress":"[=========================================\u003e         ]  27.76MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":27394560,"total":33554431},"progress":"[========================================\u003e          ]  27.39MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":27394560,"total":33554431},"progress":"[========================================\u003e          ]  27.39MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":28115456,"total":33554431},"progress":"[=========================================\u003e         ]  28.12MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":29196800,"total":33554431},"progress":"[===========================================\u003e       ]   29.2MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":29557248,"total":33554431},"progress":"[============================================\u003e      ]  29.56MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":28836352,"total":33554431},"progress":"[==========================================\u003e        ]  28.84MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":28836352,"total":33554431},"progress":"[==========================================\u003e        ]  28.84MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":29917696,"total":33554431},"progress":"[============================================\u003e      ]  29.92MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":30999040,"total":33554431},"progress":"[==============================================\u003e    ]     31MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":31359488,"total":33554431},"progress":"[==============================================\u003e    ]  31.36MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":30278144,"total":33554431},"progress":"[=============================================\u003e     ]  30.28MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":30278144,"total":33554431},"progress":"[=============================================\u003e     ]  30.28MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":31359488,"total":33554431},"progress":"[==============================================\u003e    ]  31.36MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":32801280,"total":33554431},"progress":"[================================================\u003e  ]   32.8MB/33.55MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":32440832,"total":33554431},"progress":"[================================================\u003e  ]  32.44MB/33.55MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":31719936,"total":33554431},"progress":"[===============================================\u003e   ]  31.72MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":32080384,"total":33554431},"progress":"[===============================================\u003e   ]  32.08MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":33555968,"total":33554431},"progress":"[==================================================\u003e]  33.56MB","id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":32801280,"total":33554431},"progress":"[================================================\u003e  ]   32.8MB/33.55MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":33555968,"total":33554431},"progress":"[==================================================\u003e]  33.56MB","id":"4367a98dd925"}
	{"status":"Pushing","progressDetail":{"current":33161728,"total":33554431},"progress":"[=================================================\u003e ]  33.16MB/33.55MB","id":"8407c4f3604d"}
	{"status":"Pushing","progressDetail":{"current":33555968,"total":33554431},"progress":"[==================================================\u003e]  33.56MB","id":"6a996c0ce279"}
	{"status":"Pushing","progressDetail":{"current":33522176,"total":33554431},"progress":"[=================================================\u003e ]  33.52MB/33.55MB","id":"38be7762a5d3"}
	{"status":"Pushing","progressDetail":{"current":33555968,"total":33554431},"progress":"[==================================================\u003e]  33.56MB","id":"8407c4f3604d"}
	{"status":"Pushed","progressDetail":{},"id":"d6f45f2d1604"}
	{"status":"Pushing","progressDetail":{"current":33555968,"total":33554431},"progress":"[==================================================\u003e]  33.56MB","id":"38be7762a5d3"}
	{"status":"Layer already exists","progressDetail":{},"id":"36b45d63da70"}
	{"status":"Pushed","progressDetail":{},"id":"4367a98dd925"}
	{"status":"Pushed","progressDetail":{},"id":"6a996c0ce279"}
	{"status":"Pushed","progressDetail":{},"id":"8407c4f3604d"}
	{"status":"Pushed","progressDetail":{},"id":"38be7762a5d3"}
	{"status":"tag1: digest: sha256:b85b4ed8bb804e9ebcc985bcab6dddbeb75656ed7c1186e4694d32b2b0512b35 size: 1587"}
	{"progressDetail":{},"aux":{"Tag":"tag1","Digest":"sha256:b85b4ed8bb804e9ebcc985bcab6dddbeb75656ed7c1186e4694d32b2b0512b35","Size":1587}}
`

	// 	expected := `prefix ──  ‣  [d6f45f2d1604]  Pushed
	// prefix ──  ‣  [8407c4f3604d]  Pushed
	// prefix ──  ‣  [4367a98dd925]  Pushed
	// prefix ──  ‣  [36b45d63da70]  Layer already exists
	// prefix ──  ‣  [38be7762a5d3]  Pushed
	// prefix ──  ‣  [6a996c0ce279]  Pushed

	// prefix ──  ‣  tag1: digest: sha256:b85b4ed8bb804e9ebcc985bcab6dddbeb75656ed7c1186e4694d32b2b0512b35 size: 1587
	// `

	pr, pw := io.Pipe()
	r := &DefaultResponse{
		Prefix: "prefix",
	}

	go func() {
		scanner := bufio.NewScanner(strings.NewReader(input))
		for scanner.Scan() {
			fmt.Fprintln(pw, strings.TrimSpace(scanner.Text()))
			time.Sleep(100 * time.Millisecond)
		}
		pw.Close()
	}()

	w.Reset()
	r.Write(io.Writer(&w), ioutil.NopCloser(io.Reader(pr)))

	fmt.Println(w.String())
}
