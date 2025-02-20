# PDF generation service


## Description
This project translates given input (so far only XML format) to generate PDF based on the submitted values.
It first marshals the input to the struct, generate HTML which will be out generated, and then generates the PDF based on the HTML.


## Directory explanation
- examples: It contains the examples of the input XML.
- models: It contains the struct of the input XML.
- services: It contains the business logic to generate HTML and PDF


## Requirements to run 
- Golang version 1.23.2 
- [weasyprint](https://doc.courtbouillon.org/weasyprint/stable/index.html) library to generate PDF
  - > If you are not using macOS and failed generating PDF (because of cross-compile issue), then install the library from [here](https://doc.courtbouillon.org/weasyprint/stable/first_steps.html).
  - This package was chosen for the below reason
    - It is quite popular library based on the activity logs in [GitHub](https://github.com/Kozea/WeasyPrint)
    - It is easy to use compared to using [maroto](https://github.com/johnfercher/maroto) library to generate PDF for dynamic content
    - [wkhtmltopdf](https://github.com/SebastiaanKlippert/go-wkhtmltopdf) was considered, but it is no longer maintained


## How to run
1. Fulfill all the requirements above
2. place the XML file in the root directory with the name of `input.xml`.
3. Run the command below
    ```shell
    go run main.go
    ```
4. PDF and HTML will be generated in the root directory


## Service specification
This service is available to parse the below conditions
- Inside the `Form` element
    - It has to have more than 1 `Section` and `Field` elements
- Inside the `Section` element
    - `Title` and `Contents` are mandatory, but can’t have more than 2 for each.
    - `Contents` has to have more than 1 `Section` or `Field` elements  
    - This element can takes only “Optional” and "Name" attribute.
- Inside the `Field` element
    - `Caption` is mandatory, but it takes the first `Caption` element if there is more than one.
    - Validation will be applied to the `Label` elements based on the “FieldType” attribute. Below is the existing validations.
      - TextBox: Inside the "Type" attribute, it should have minLength, maxLength and rows as HTML textarea attribute values. <br> Example schema: `Text([0,200],Lines:4)`  
      - File: No validation
      - Select: Inside the "Type" attribute, it should list the available enumeration value where `<Label>` element can use. <br> If `<Label>` element uses the value in Name attribute which is not in the list, it will be considered as invalid. <br> Example schema: `Enumeration([A,B,C])`


## Future extension
- Support more input field types 
  - Steps:
    - Add the desired field types for the struct in models/models.go
    - Add the validation and business logic for the desired field types in convertFieldToHTML() in services/services.go
- Support more form elements in addition to fields/inputs and sections, such as comment boxes.
  - Steps:
    - Add the desired form elements for the struct in models/models.go
    - Add the business logic in GenerateHTML() in services/services.go as needed
- Support different input form schema data structures (e.g., FormIO’s JSON).
  - Steps:
    - Try to fit the existing Form struct in models/models.go
    - Add the business logic for parsing in parseJSON() in services/parser.go 