require 'sinatra'
get '/' do
  address = ENV["NAME"]

  if address
    "Hello #{address}, Welcome to Waypoint!"
  else
    'Welcome to Waypoint!'
  end
end
