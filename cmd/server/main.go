package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/admein-advertising/admein-vast-generator/validator"
	"github.com/admein-advertising/admein-vast-generator/vast"
)

const vastMedia1 = `https://admein.io/api/stream?id=52dab5e4-0-480p`

func main() {
	http.HandleFunc("/vast", withCORS(vastHandler))
	http.HandleFunc("/vast/validate", withCORS(vastValidateHandler))
	http.HandleFunc("/vast/example1", withCORS(vastExample1Handler))
	http.HandleFunc("/vast/example2", withCORS(vastExample2Handler))
	http.HandleFunc("/vast/example3", withCORS(vastExample3Handler))
	http.HandleFunc("/vast/example4", withCORS(vastExample4Handler))
	http.HandleFunc("/vast/example5", withCORS(vastExample5Handler))
	http.HandleFunc("/", withCORS(homeHandler))
	fmt.Println("Server started at http://localhost:3780")
	log.Fatal(http.ListenAndServe(":3780", nil))
}

// Example 1: Simple VAST with one ad
// Keeping the first definition of vastExample1Handler
func vastExample1Handler(w http.ResponseWriter, r *http.Request) {
	// Create a new VAST document
	v := vast.New()

	// Add an inline ad
	duration, err := vast.NewDuration("00:00:05")
	if err != nil {
		log.Fatal(err)
	}
	v.Ad = append(v.Ad, vast.Ad{
		ID:       "1",
		Sequence: 1,
		InLine: &vast.InLine{
			AdTitle: "AdMeIn VAST Example",
			Creatives: vast.InLineCreatives{
				Creative: []vast.InLineCreative{
					{
						Creative: vast.Creative{Sequence: 1},
						Linear: &vast.LinearInLine{
							Duration: duration,
							MediaFiles: vast.MediaFiles{
								MediaFile: []vast.MediaFile{
									{
										Value:    vastMedia1,
										Delivery: vast.ProgressiveDelivery,
										Type:     "video/mp4",
										Width:    640,
										Height:   360,
									},
								},
							},
						},
					},
				},
			},
		},
	})
	xmlBytes, err := v.Bytes()
	if err != nil {
		http.Error(w, "Failed to marshal VAST", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

// Example 2: VAST with two ads in a pod
// Keeping the first definition of vastExample2Handler
func vastExample2Handler(w http.ResponseWriter, r *http.Request) {
	v := vast.New()
	v.Ad = append(v.Ad,
		vast.Ad{
			ID:       "1",
			Sequence: 1,
			InLine: &vast.InLine{
				AdTitle: "Example 2 Ad 1",
				Creatives: vast.InLineCreatives{
					Creative: []vast.InLineCreative{
						{
							Creative: vast.Creative{Sequence: 1},
							Linear: &vast.LinearInLine{
								Duration: "00:00:08",
								MediaFiles: vast.MediaFiles{
									MediaFile: []vast.MediaFile{
										{
											Value:    vastMedia1,
											Delivery: vast.ProgressiveDelivery,
											Type:     "video/mp4",
											Width:    640,
											Height:   360,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		vast.Ad{
			ID:       "2",
			Sequence: 2,
			InLine: &vast.InLine{
				AdTitle: "Example 2 Ad 2",
				Creatives: vast.InLineCreatives{
					Creative: []vast.InLineCreative{
						{
							Creative: vast.Creative{Sequence: 2},
							Linear: &vast.LinearInLine{
								Duration: "00:00:06",
								MediaFiles: vast.MediaFiles{
									MediaFile: []vast.MediaFile{
										{
											Value:    vastMedia1,
											Delivery: vast.ProgressiveDelivery,
											Type:     "video/mp4",
											Width:    640,
											Height:   360,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
	xmlBytes, err := v.Bytes()
	if err != nil {
		http.Error(w, "Failed to marshal VAST", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

// Example 3: VAST ad pod with three ads
// Keeping the first definition of vastExample3Handler
func vastExample3Handler(w http.ResponseWriter, r *http.Request) {
	v := vast.New()
	v.Ad = append(v.Ad,
		vast.Ad{
			ID:       "1",
			Sequence: 1,
			InLine: &vast.InLine{
				AdTitle: "Example 3 Ad 1",
				Creatives: vast.InLineCreatives{
					Creative: []vast.InLineCreative{
						{
							Creative: vast.Creative{Sequence: 1},
							Linear: &vast.LinearInLine{
								Duration: "00:00:10",
								MediaFiles: vast.MediaFiles{
									MediaFile: []vast.MediaFile{
										{
											Value:    vastMedia1,
											Delivery: vast.ProgressiveDelivery,
											Type:     "video/mp4",
											Width:    640,
											Height:   360,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		vast.Ad{
			ID:       "2",
			Sequence: 2,
			InLine: &vast.InLine{
				AdTitle: "Example 3 Ad 2",
				Creatives: vast.InLineCreatives{
					Creative: []vast.InLineCreative{
						{
							Creative: vast.Creative{Sequence: 2},
							Linear: &vast.LinearInLine{
								Duration: "00:00:12",
								MediaFiles: vast.MediaFiles{
									MediaFile: []vast.MediaFile{
										{
											Value:    vastMedia1,
											Delivery: vast.ProgressiveDelivery,
											Type:     "video/mp4",
											Width:    640,
											Height:   360,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		vast.Ad{
			ID:       "3",
			Sequence: 3,
			InLine: &vast.InLine{
				AdTitle: "Example 3 Ad 3",
				Creatives: vast.InLineCreatives{
					Creative: []vast.InLineCreative{
						{
							Creative: vast.Creative{Sequence: 3},
							Linear: &vast.LinearInLine{
								Duration: "00:00:08",
								MediaFiles: vast.MediaFiles{
									MediaFile: []vast.MediaFile{
										{
											Value:    vastMedia1,
											Delivery: vast.ProgressiveDelivery,
											Type:     "video/mp4",
											Width:    640,
											Height:   360,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
	xmlBytes, err := v.Bytes()
	if err != nil {
		http.Error(w, "Failed to marshal VAST", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

func vastExample4Handler(w http.ResponseWriter, r *http.Request) {
	v := vast.New()
	v.Ad = append(v.Ad, vast.Ad{
		ID:       "1",
		Sequence: 1,
		Wrapper: &vast.Wrapper{
			AdVerifications: &vast.AdVerifications{
				Verification: []vast.Verification{
					{
						JavaScriptResource: []vast.JavaScriptResource{
							{Value: "https://example.com/verification.js"},
						},
					},
				},
			},
			AdDefinition: vast.AdDefinition{
				AdSystem: vast.AdSystem{Value: "Example Wrapper System"},
				Impression: []vast.Impression{
					{Value: "https://example.com/impression"},
				},
			},
			VASTAdTagURI: vast.CData{
				Value: "https://example.vudoo.io/backmagic/ads/vast/4.0/88995844754/ctv?tag_id=89598758273",
			},
			Creatives: &vast.Creatives{
				Creative: []vast.WrapperCreative{
					{
						ID: "1",
						Linear: &vast.LinearWrapper{
							Linear: vast.Linear{
								TrackingEvents: &vast.TrackingEvents{
									Tracking: []vast.Tracking{
										{Value: "https://example.com/track"},
									},
								},
							},
						},
					},
				},
			},
			FollowAdditionalWrappers: vast.NumericBool(false),
			AllowMultipleAds:         vast.NumericBool(false),
			FallbackOnNoAd:           vast.NumericBool(false),
		},
	})
	xmlBytes, err := v.Bytes()
	if err != nil {
		http.Error(w, "Failed to marshal VAST Wrapper", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

// Read from URL example
func vastExample5Handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://admein-advertising.github.io/vast-tag-xml-examples/vast-4.1/universal-ad-Id-vast-4-1-sample.xml")
	if err != nil {
		http.Error(w, "Failed to fetch VAST from URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	example, err := vast.Read(resp.Body)
	if err != nil {
		http.Error(w, "Failed to parse VAST", http.StatusInternalServerError)
		return
	}

	xmlBytes, err := example.Bytes()
	if err != nil {
		http.Error(w, "Failed to marshal VAST", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

// withCORS is a middleware to add CORS headers
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<h1>VAST Server examples</h1><p>Use <a href='/vast' target="_blank">/vast</a> to get a VAST response or see the below examples.</p>`)
	fmt.Fprintf(w, `<p>Examples:</p><ul>
		<li><a href='/vast/example1' target='_blank'>/vast/example1</a></li>
		<li><a href='/vast/example2' target='_blank'>/vast/example2</a></li>
		<li><a href='/vast/example3' target='_blank'>/vast/example3</a></li>
		<li><a href='/vast/example4' target='_blank'>/vast/example4</a></li>
		<li><a href='/vast/example5' target='_blank'>/vast/example5</a></li>
		<li><strong>POST</strong> raw XML to <code>/vast/validate</code> to receive a JSON validation report.</li>
	</ul>`)
}

func vastHandler(w http.ResponseWriter, r *http.Request) {
	v := vast.New()

	// Example: Add 1 Ad with 1 Creative to the VAST response
	v.Ad = append(v.Ad, vast.Ad{
		ID:       "1",
		Sequence: 1,
		InLine: &vast.InLine{
			AdTitle: "Sample Ad",
			Creatives: vast.InLineCreatives{
				Creative: []vast.InLineCreative{
					{
						Creative: vast.Creative{Sequence: 1},
						Linear: &vast.LinearInLine{
							Duration: "00:00:10",
							MediaFiles: vast.MediaFiles{
								MediaFile: []vast.MediaFile{
									{
										Value:    vastMedia1,
										Delivery: vast.ProgressiveDelivery,
										Type:     "video/mp4",
										Width:    640,
										Height:   360,
									},
								},
							},
						},
					},
				},
			},
		},
	})

	xmlBytes, err := v.Bytes()
	if err != nil {
		http.Error(w, "Failed to marshal VAST", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

func vastValidateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	result, err := validator.Validate(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode validation result: %v", err)
	}
}
