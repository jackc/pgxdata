require "base64"
require "erb"
require "fileutils"
require "rake/clean"

CLEAN.include "templates.go"

file "templates.go" => ["templates.go.erb", *FileList["templates/*"]] do |t|
  templates = Dir.glob("templates/*").each_with_object({}) do |path, h|
    h[File.basename(path)] = File.read(path)
  end
  erb = ERB.new(File.read("templates.go.erb"))
  File.write("templates.go", "// Do not edit. Automatically generated from templates.go.erb\n" + erb.result(binding))
  sh "gofmt -w templates.go"
end

desc "Run go tests"
task :test => ["templates.go"] do
  sh "go test ./..."
end

task :default => [:test]
