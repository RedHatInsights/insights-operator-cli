pushd ..
goplantuml -recursive . > class_diagram.uml
java -jar ~/tools/plantuml.jar class_diagram.uml
mv class_diagram.uml doc/
mv class_diagram.png doc/
popd
