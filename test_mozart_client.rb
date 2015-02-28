require_relative './mozart_client'

correct_task = {
  in: '10s',
  do: 'http_request',
  timeout: '15s',
  params: {
    url: 'http://localhost:4000/task_executed',
    verb: 'POST',
    header: { 'X-Auth-Token': 'MY_TOKEN' },
    json_body: { foo: 'bar', number: 123 }
  }
}

another_correct_task = {
  in: '3s',
  do: 'write_file',
  params: {
    filename: '/tmp/foobar',
    filemode: 0640,
    json_body: { foo: 'bar', number: 123 }
  }
}

mozart_client = MozartClient.new
puts "Count: #{mozart_client.count_tasks}"

resp_data = mozart_client.schedule_task(correct_task)
puts "Task created with uuid: #{resp_data['task_uuid']}"

resp_data = mozart_client.schedule_task(another_correct_task)
puts "Task created with uuid: #{resp_data['task_uuid']}"

puts "Count: #{mozart_client.count_tasks}"

invalid_task = correct_task
invalid_task[:in] = 'invalid time value'
invalid_task[:do] = 'foo'
begin
  mozart_client.schedule_task(invalid_task)
rescue ArgumentError => e
  puts e
end

mozart_client.each_scheduled_task do |task|
  puts "Scheduled task: #{task['uuid']}"
  # mozart_client.unschedule_task(task['uuid'])
end

Thread.abort_on_exception = true
thread = Thread.new do
  class App < Sinatra::Base
    post '/task_executed' do
      puts '** Web server receive: Task executed!'
    end
  end
  App.run!(host: 'localhost', port: 4000)
end

thread.join

