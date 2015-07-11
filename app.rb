require 'sinatra'
require 'json'
require 'pp'

post '/scheduled_tasks' do
  pp env
  # data = JSON.parse(request.body.read)
  # pp data
  "OK"
end

