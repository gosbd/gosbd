# **gosbd: Sentence Boundary Disambiguation (SBD) Library for Go**

<img align="right" width="320" src="/artifacts/sbd-gopher.png" alt="gosbd-logo" title="dsbd-logo" />

[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gosbd/gosbd)

gosbd is a library for segmenting text into sentences inspired by [pySBD](https://github.com/nipunsadvilkar/pySBD) and [pragmatic_segmenter](https://github.com/diasks2/pragmatic_segmenter). It is rule-based and works out-of-the-box.

## Features

- **Sentence Segmentation**: Efficiently breaks down a block of text into individual sentences.
- **Lightweight and Easy Integration**: Designed to be lightweight and easy to integrate into existing Go projects.
- **High Accuracy**: Offers high accuracy in sentence segmentation. For more details, see [pySBD](https://github.com/nipunsadvilkar/pySBD).
- **Non-Destructive Segmentation**: Segments text into sentences without altering the original content.
- **Language-Specific Configuration**: Adaptable to handle punctuation rules specific to different languages.
- **Text Cleaning**: Equipped with features to manage and clean noisy text, including:
    - Handling irregular newline characters and spacing
    - Processing Tables of Contents
    - Recognizing and managing URLs and HTML tags
    - Dealing with sentences that are delimited without any space

_Note: Text Cleaning feature is to be implemented. Contributions are greatly welcomed._

## Installation

To install gosbd, you can use `go get`:

```sh
go get github.com/gosbd/gosbd
```

## Usage

Here's a basic example of how to use gosbd:

```go
package main

import (
    "fmt"
    "github.com/gosbd/gosbd"
)

// This example segments a text string into individual sentences.
func main() {
    segmenter := gosbd.NewSegmenter("en")
    text := "This is a sentence. And this is another one."
    sentences := segmenter.Segment(text)
    for _, sentence := range sentences {
        fmt.Println(sentence)
    }
}
```

## Roadmap

- [ ] Add support for more languages.
- [ ] Implement text cleaner.
- [ ] Add benchmark test.
- [ ] Setup GitHub Action for testing.
- [ ] Setup Codecov for monitoring test coverage.
- [ ] Add Online Playground.

## Language Support Roadmap

The following table outlines our current language support. We're actively seeking contributions to expand this list. If you're interested in contributing, consider helping us add support for a language, whether it's listed below or not. Your expertise in a language not listed here could be a valuable addition to our project.

| Language   | ISO Code | Supported |
| ---------- | -------- |-----------|
| Amharic    | am       | Planned   |
| Arabic     | ar       | Planned   |
| Armenian   | hy       | Planned   |
| Bulgarian  | bg       | Planned   |
| Burmese    | my       | Planned   |
| Chinese    | zh       | Yes       |
| Danish     | da       | Planned   |
| Deutsch    | de       | Planned   |
| Dutch      | nl       | Planned   |
| English    | en       | Yes       |
| French     | fr       | Planned   |
| Greek      | el       | Planned   |
| Hindi      | hi       | Planned   |
| Italian    | it       | Planned   |
| Japanese   | ja       | Yes       |
| Kazakh     | kk       | Planned   |
| Marathi    | mr       | Planned   |
| Persian    | fa       | Planned   |
| Polish     | pl       | Planned   |
| Russian    | ru       | Yes       |
| Slovak     | sk       | Planned   |
| Spanish    | es       | Planned   |
| Urdu       | ur       | Planned   |

We welcome contributions that help us add support for these languages. Please feel free to submit a Pull Request with your contributions.

## Motivation

Sentence Boundary Disambiguation (SBD) is a step in the preprocessing pipeline of Natural Language Processing (NLP). It involves the task of correctly identifying the boundaries of sentences within a block of text, which is fundamental for subsequent tasks such as Machine Translation, Named Entity Recognition, and Coreference Resolution.

The importance of accurate sentence segmentation has been amplified with the widespread use of Large Language Models (LLMs). While the libraries [pragmatic_segmenter](https://github.com/diasks2/pragmatic_segmenter) and [pySBD](https://github.com/nipunsadvilkar/pySBD) are known for their high accuracy and efficiency, there are no equivalent libraries available in Go.

This library seeks to bridge this gap by offering a rule-based sentence segmentation solution in Go. The ability to perform sentence segmentation in Go offers significant advantages, especially when using LLMs via an API. For instance, it allows tasks such as Retrieval Augmented Generation (RAG) and summarization to be completed entirely within Go. This not only streamlines the development process, but it can also lead to faster execution times due to Go's performance characteristics.

## Acknowledgement

This library builds upon the excellent foundations laid by [pySBD](https://github.com/nipunsadvilkar/pySBD) and [pragmatic_segmenter](https://github.com/diasks2/pragmatic_segmenter).

## Contributing

Contributions are greatly appreciated and crucial for this project! Here are a few ways you can contribute:

- **Add new tests and rules**: Improve the accuracy of sentence segmentation by adding new tests and rules.
- **Add support for a new language**: Help expand the reach of this library by adding support for new languages.
- **Port features**: Help improve this library by porting features that are supported in pySBD and pragmatic_segmenter.

Please feel free to submit a Pull Request with your contributions.

## License

This project is licensed under the MIT License.
