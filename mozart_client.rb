require 'net/http'
require 'json'
require 'sinatra'

class MozartClient

  DEFAULT_ENDPOINT = 'http://localhost:1357'

  InvalidRespStatusCode = Class.new(StandardError)

  attr_reader :endpoint

  def initialize(endpoint=DEFAULT_ENDPOINT)
    @endpoint = endpoint
  end

  def count_tasks
    path = '/tasks/count'
    resp = get(uri_with(path))
    raise InvalidRespStatusCode.new(resp.code) if resp.code != '200'
    data = parse_json(resp.body)
    data['count']
  end

  def schedule_task(data)
    path = '/tasks/schedule'
    resp = post(uri_with(path), data)
    if resp.code == '400'
      data = parse_json(resp.body)
      raise ArgumentError.new( data['errors'].join(', ') )
    end
    raise InvalidRespStatusCode.new(resp.code) if resp.code != '202'
    parse_json(resp.body)
  end

  def unschedule_task(uuid)
    path = "/tasks/unschedule/#{uuid}"
    delete(uri_with(path))
  end

  def each_scheduled_task
    path = "/tasks"
    resp = get(uri_with(path))
    raise InvalidRespStatusCode.new(resp.code) if resp.code != '202'
    parse_json(resp.body).each do |task|
      yield task
    end
  end

  private

  def delete(uri)
    delete = Net::HTTP::Delete.new(uri.request_uri)
    resp = http(uri).request(delete)
  end

  def post(uri, data)
    post = Net::HTTP::Post.new(uri.request_uri)
    post.body = data.to_json
    http(uri).request(post)
  end

  def get(uri)
    get = Net::HTTP::Get.new(uri.request_uri)
    resp = http(uri).request(get)
  end

  def http(uri)
    http = Net::HTTP.new(uri.host, uri.port)
  end

  def uri_with(path)
    URI.parse("#{endpoint}#{path}")
  end

  def parse_json(raw_json)
    JSON.parse(raw_json)
  end

end

