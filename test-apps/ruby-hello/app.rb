require 'sinatra'

def handler(event:, context:)
  { event: JSON.generate(event), context: JSON.generate(context.inspect) }
end

set :bind, '0.0.0.0'
set :port, 8080

get '/' do
  'Hello World!'
end

