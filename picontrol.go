package engine

/*
#cgo LDFLAGS:-Llib -lwiringPi

#include <stdio.h>
#include <wiringPi.h>

#define PWM0_0 1
#define PWM0_1 26
#define PWM1_0 23
#define PWM1_1 24
#define CAR_TYPE_PIN 4
#define MODE PWM_MODE_MS
#define BLADE_SWITCH 5

//int lastSpeed=0;
//int lastDirection=0;

int wiringInit(){
  wiringPiSetup();
  pinMode(CAR_TYPE_PIN,INPUT);
  pinMode(BLADE_SWITCH,OUTPUT);
  pinMode(PWM0_1,OUTPUT);
  pinMode(PWM0_0,PWM_OUTPUT);

  pinMode(PWM1_1,OUTPUT);
  pinMode(PWM1_0,PWM_OUTPUT);

  pwmSetMode(MODE);
  pwmWrite(PWM1_0,0);
  pwmWrite(PWM0_0,0);
  digitalWrite(PWM0_1,0);
  digitalWrite(PWM1_1,0);
  return digitalRead(CAR_TYPE_PIN);
}

 // 两轮差速控制小车的速度控制
void speedControl0(int lastSpeed,int speed){
	if(speed >=0){
	 if(speed>16) speed=16;
	   if(lastSpeed >=0)
	   {
	     pwmWrite(PWM0_0,64*speed);
	     pwmWrite(PWM1_0,64*speed);
	   }else{  //speed>=0 and lastSpeed<0
			pinMode(PWM0_0,PWM_OUTPUT);
			pinMode(PWM0_1,OUTPUT);
			pinMode(PWM1_0,PWM_OUTPUT);
			pinMode(PWM1_1,OUTPUT);
			digitalWrite(PWM0_1,0);
			digitalWrite(PWM1_1,0);
			pwmSetMode(MODE);
			pwmWrite(PWM0_0,64*speed);
			pwmWrite(PWM1_0,64*speed);
	    }
	 }else{  //speed<0

	  if(speed<-16) speed=-16;

	  if(lastSpeed <0) //speed<0 and lastSpeed<0
	   {
	       pwmWrite(PWM0_1,-64*speed);
	       pwmWrite(PWM1_1,-64*speed);
	   }else{ //speed<0 and lastSpeed>=0
			pinMode(PWM0_1,PWM_OUTPUT);
			pinMode(PWM0_0,OUTPUT);
			pinMode(PWM1_1,PWM_OUTPUT);
			pinMode(PWM1_0,OUTPUT);
			digitalWrite(PWM0_0,0);
			digitalWrite(PWM1_0,0);
			pwmSetMode(MODE);
			pwmWrite(PWM0_1,-64*speed);
			pwmWrite(PWM1_1,-64*speed);
	    }
	 }//else
	//lastSpeed=speed;
	}

	 // 两轮差速控制小车的方向控制
void directionControl0(int lastSpeed,int direction){
    //int level=lastSpeed+8;. case 0:
    if(direction==0){
      speedControl0(lastSpeed,lastSpeed);
      return;
      }
    direction=direction/2;//减慢速度
    int tempSpeed=0;
    switch (lastSpeed){
      case -16:
      case -14:
      case -12:
      case -10:
         if(direction >=0 ){  //turn right
	  tempSpeed=lastSpeed+direction;
	  if(tempSpeed>0) tempSpeed=0;
	  pwmWrite(PWM0_0,-64*tempSpeed);
	 }else{  //turn left
	  tempSpeed=lastSpeed-direction;
	  if(tempSpeed<-16) tempSpeed=-16;
	  pwmWrite(PWM1_0,-64*tempSpeed);
	 }
       break;
      case -8:
      case -6:
      case -4:
      case -2:
         if(direction >=0 ){  //turn right
	  tempSpeed=lastSpeed-direction;
	  if(tempSpeed<-16) tempSpeed=-16;
	  pwmWrite(PWM1_0,-64*tempSpeed);
	 }else{  //turn left
	  tempSpeed=lastSpeed+direction;
	  if(tempSpeed>0) tempSpeed=0;
	  pwmWrite(PWM0_0,-64*tempSpeed);
	 }
        break;
	  case 0:
      case 2:
      case 4:
      case 6:
      case 8:
        if(direction >=0 ){  //turn right
	  tempSpeed=lastSpeed+direction;
	  if(tempSpeed>16) tempSpeed=16;
	  pwmWrite(PWM1_0,64*tempSpeed);
	 }else{  //turn left
	  tempSpeed=lastSpeed-direction;
	  if(tempSpeed>16) tempSpeed=16;
	  pwmWrite(PWM0_0,64*tempSpeed);
	 }
       break;

      case 10:
      case 12:
      case 14:
      case 16:
     if(direction >=0 ){  //turn right
	  tempSpeed=lastSpeed-direction;
	  if(tempSpeed<0) tempSpeed=0;
	  pwmWrite(PWM0_0,64*tempSpeed);
	 }else{  //turn left,direction<0
	  tempSpeed=lastSpeed+direction;
	  if(tempSpeed<0) tempSpeed=0;
	  pwmWrite(PWM1_0,64*tempSpeed);
	 }
       break;
      default:break;
    }
  }

//四轮阿克曼转向架控制小车的速度控制
void speedControl1(int lastSpeed,int speed){
	if(speed >=0){
	 if(speed>16) speed=16;

	   if(lastSpeed >=0)
	   {
	     pwmWrite(PWM0_0,64*speed);

	   }else{  //speed>=0 and lastSpeed<0
			pinMode(PWM0_0,PWM_OUTPUT);
			pinMode(PWM0_1,OUTPUT);
			digitalWrite(PWM0_1,0);
			pwmSetMode(MODE);
			pwmWrite(PWM0_0,64*speed);
	    }
	 }else{  //speed<0

	  if(speed<-16) speed=-16;

	  if(lastSpeed <0) //speed<0 and lastSpeed<0
	   {
	     pwmWrite(PWM0_1,-64*speed);
	   }else{ //speed<0 and lastSpeed>=0
			pinMode(PWM0_1,PWM_OUTPUT);
			pinMode(PWM0_0,OUTPUT);
			digitalWrite(PWM0_0,0);
			pwmSetMode(MODE);
			pwmWrite(PWM0_1,-64*speed);
	    }
	 }//else
	//lastSpeed=speed;
	}

	//四轮阿克曼转向架控制小车的方向控制
void directionControl1(int lastDirection,int direction){
   if(direction >=0){
	 if(direction>16) direction=16;

	   if(lastDirection >=0)
	   {
	     pwmWrite(PWM1_0,64*direction);

	   }else{  //speed>=0 and lastSpeed<0
			pinMode(PWM1_0,PWM_OUTPUT);
			pinMode(PWM1_1,OUTPUT);
			digitalWrite(PWM1_1,0);
			pwmSetMode(MODE);
			pwmWrite(PWM1_0,64*direction);
	    }
	 }else{  //speed<0

	  if(direction<-16) direction=-16;

	  if(lastDirection <0) //speed<0 and lastSpeed<0
	   {
	     pwmWrite(PWM1_1,-64*direction);
	   }else{ //speed<0 and lastSpeed>=0
			pinMode(PWM1_1,PWM_OUTPUT);
			pinMode(PWM1_0,OUTPUT);
			digitalWrite(PWM1_0,0);
			pwmSetMode(MODE);
			pwmWrite(PWM1_1,-64*direction);
	    }
	 }//else
	//lastSpeed=speed;
}

//刀片开关控制
void bladeSwitch(int bladeOn){
	 digitalWrite(BLADE_SWITCH,bladeOn);
	// printf("blade:%d\n",bladeOn);
}
*/
import "C"

type PiControl struct {
	lastSpeed     int
	lastDirection int
	carType       C.int
}

func Init() *PiControl {
	_carType := C.wiringInit()
	//fmt.Printf("_carType: %d\n", _carType)
	return &PiControl{
		lastSpeed: 0,
		carType:   _carType,
	}
}

func (pi *PiControl) SpeedControl(speed int) error {
	if pi.carType == 0 { // 两轮差速控制小车的速度控制
		C.speedControl0(C.int(pi.lastSpeed), C.int(speed))
		//	fmt.Printf("speed type: %d\n", pi.carType)
	} else { //pi.carType == 1 四轮阿克曼转向架控制小车的速度控制
		C.speedControl1(C.int(pi.lastSpeed), C.int(speed))
	}
	pi.lastSpeed = speed
	return nil
}

// 两轮差速控制小车的方向控制
func (pi *PiControl) DirectionControl(direction int) error {
	if pi.carType == 0 { // 两轮差速控制小车的方向控制
		C.directionControl0(C.int(pi.lastSpeed), C.int(direction))
	} else { //pi.carType == 1 四轮阿克曼转向架控制小车的方向控制
		C.directionControl1(C.int(pi.lastDirection), C.int(direction))
		//	fmt.Printf("dir type: %d\n", pi.carType)
	}
	pi.lastDirection = direction
	return nil
}

// 刀片开关
func (pi *PiControl) BladeSwitch(bladeOn int) error {
	C.bladeSwitch(C.int(bladeOn))
	//fmt.Printf("blade: %d\n", bladeOn)
	return nil
}
