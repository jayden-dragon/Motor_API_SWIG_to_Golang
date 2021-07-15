%module(director="1") Motor

%include <typemaps.i>

%{
#include <fcntl.h>
#include <termios.h>
#include <unistd.h>

#include <stdlib.h>
#include <stdio.h>

#include "dynamixel_sdk.h"
#include "port_handler.h"
#include "packet_handler.h"
%}

%feature("director") PortHandlerLinux;
%feature("director") Protocol2PacketHandler;


%apply int {uint8_t}
%apply int {uint16_t}
%apply double {float}

%include "dynamixel_sdk.h"
%include "port_handler.h"
%include "packet_handler.h" 
