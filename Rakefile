require "base64"
require "erb"
require "fileutils"
require "rake/clean"

CLEAN.include "templates.go"
CLEAN.include FileList["test/data/*.go"].exclude("test/data/*_test.go")
CLOBBER.include "build/*"

file "templates.go" => ["templates.go.erb", *FileList["templates/*"]] do |t|
  templates = Dir.glob("templates/*").sort.each_with_object({}) do |path, h|
    h[File.basename(path)] = File.read(path)
  end
  erb = ERB.new(File.read("templates.go.erb"))
  File.write("templates.go", "// Do not edit. Automatically generated from templates.go.erb\n" + erb.result(binding))
  sh "gofmt -w templates.go"
end

file "test/data/db.go" => ["build/pgxdata", "test/data/config.toml"] do
  sh "cd test/data && ../../build/pgxdata generate && gofmt -w *.go"
end

file "build/pgxdata" => FileList["*.go"] do
  Dir.mkdir("build") unless Dir.exists?("build")
  sh "go build -o build/pgxdata"
end

desc "Run go tests"
task :test => FileList["templates.go", "test/data/db.go"] do
  sh "go test ./..."
end

desc "Build pgxdata"
task :build => "build/pgxdata"

task :default => [:test]
