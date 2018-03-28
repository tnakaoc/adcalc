package main
import "fmt"
import "os"
import "bufio"
import "strconv"
func main(){
	if len(os.Args)==1 {
		fmt.Println("usage : ",os.Args[0]," [filename] ([normal dose])")
		return
	}
	file,err:=os.Open(os.Args[1])
	if err!=nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	limitNorm,ppmNorm,ratTB:=func()(float64,float64,float64){
		if len(os.Args)==3 {
			v,e:=strconv.ParseFloat(os.Args[2],64)
			if e==nil {
				return v,25.0,3.5
			}
		}else if len(os.Args)>3 {
			ln:=13.0
			pn:=25.0
			rt:=3.5
			for i:=2;i<len(os.Args);i++ {
				arg:=os.Args[i]
				switch arg[:2] {
					case "-D":
						v,e:=strconv.ParseFloat(arg[2:],64)
						if e==nil {
							ln=v
						}
					case "-N":
						v,e:=strconv.ParseFloat(arg[2:],64)
						if e==nil {
							pn=v
						}
					case "-T":
						v,e:=strconv.ParseFloat(arg[2:],64)
						if e==nil {
							rt=v/pn
						}
					case "-R":
						v,e:=strconv.ParseFloat(arg[2:],64)
						if e==nil {
							rt=v
						}
				}
			}
			return ln,pn,rt
		}
		return 13.0,25.0,3.5
	}()
	val,err:=func()([][5]float64,error){
		scanner:=bufio.NewScanner(file)
		ret:=make([][5]float64,0)
		for scanner.Scan() {
			line:=scanner.Text()
			if line[0]=='#' { continue }
			var tmp [5]float64
			n,e:=fmt.Sscanf(line,"%v %v %v %v %v",&tmp[0],&tmp[1],&tmp[2],&tmp[3],&tmp[4])
			if e!=nil { continue }
			if n==5 {
				ret=append(ret,tmp)
			}
		}
		if len(ret)==0 {
			return ret,fmt.Errorf("invalid filetype.\n")
		}
		return ret,nil
	}()
	if err!=nil {
		fmt.Println(err)
		return
	}
	ratSB:=1.2
	cbSkin:=2.5
	cbNorm:=1.35
	cbTumor:=3.8
	limitSkin:=15.0
	max:=func(a float64,b float64)float64{
		if a>b { return a }
		return b
	}
	timeLimit,doseLimit,doseTotal:=func()(float64,float64,float64){
		timeSkin:=(val[0][2]+val[0][3]*ratSB*ppmNorm*cbSkin+val[0][4])/limitSkin
		timeNorm,totalDose:=func()(float64,float64){
			ret:=0.0
			sum:=0.0
			for _,v:=range(val) {
				tmp:=v[2]+v[3]*ppmNorm*cbNorm+v[4]
				sum+=tmp
				ret=max(tmp,ret)
			}
			return ret,sum
		}()
		timeNorm/=limitNorm
		if timeNorm>timeSkin {
			return timeNorm,limitNorm,totalDose/timeNorm
		}
		return timeSkin,limitSkin,totalDose/timeSkin
	}()
	ad,ad30,peakDose:=func()(float64,float64,float64){
		prev:=0.0
		maxd:=0.0
		ret:=[2]float64{0.0,0.0}
		for i,v:=range(val) {
			if i==0 { continue }
			tmp:=(v[2]+v[3]*ratTB*ppmNorm*cbTumor+v[4])/timeLimit
			maxd=max(tmp,maxd)
			if tmp-prev==0.0 { continue }
			if (doseLimit-tmp)*(doseLimit-prev)<=0.0 {
				ret[0]=v[0]-(tmp-doseLimit)/(tmp-prev)*(v[0]-val[i-1][0])
			}
			if prev>30.0 && tmp<=30.0 {
				ret[1]=v[0]-(tmp-30.0)/(tmp-prev)*(v[0]-val[i-1][0])
			}
			prev=tmp
		}
		return ret[0],ret[1],maxd
	}()
	fmt.Printf("%s\t%e\t%e\t%e\t%e\t%e\n",os.Args[1],doseLimit,ad,ad30,peakDose,doseTotal)
	return
}
