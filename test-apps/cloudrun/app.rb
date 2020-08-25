require 'sinatra'
get '/' do
  address = ENV["NAME"]

  if address
    "<H1>Hello #{address}, Welcome to Waypoint!</H1><br/><br/><iframe width='560' height='315' src='https://www.youtube.com/embed/OK1GDkvFdL4' frameborder='0' allow='accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture' allowfullscreen></iframe>"
  else
    'Welcome to Waypoint!'
  end
end
