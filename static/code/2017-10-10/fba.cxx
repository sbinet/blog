// g++ -o forward_backward_asymmetry forward_backward_asymmetry.cpp -I$ROOTSYS/include -L$ROOTSYS/lib  -I$ROOTSYS/include -L$ROOTSYS/lib -lCore -lCint -lRIO -lNet -lHist -lGraf -lGraf3d -lGpad -lTree -lRint -lPostscript -lMatrix -lPhysics -lMathCore -lThread -pthread -lm -ldl -rdynamic -lMinuit
#include <cstdio> 
#include <cmath> 
#include "TFile.h"
#include "TH1.h"
#include "TMinuit.h"

#define MAXEVENTS 1000

/* globale Variablen */
int nevents;
double costh[MAXEVENTS];


void results(TMinuit *minuit)
{
// helper function for nice formatting 

#define MAXPAR 50        

   double fmin,fedm;   
   double errdef;    
   int    nparv,nparx; 
   int    fstat;      

   TString pname;         
   double pvalue,perror; 
   double plbound,pubound;
   int    pvari;        
   double eplus,eminus;  
   double gcorr;         

   double emat[MAXPAR][MAXPAR]; 
   double kmat[MAXPAR][MAXPAR]; 

   int i,j;                                 
 
   minuit->mnstat(fmin,fedm,errdef,nparv,nparx,fstat);

   printf("\n\n");
   printf("Results of MINUIT minimisation\n");
   printf("-------------------------------------\n\n");
   
   printf(" Minimal function value:              %8.3lf  \n",fmin);
   printf(" Estimated difference to true minimum: %11.3le \n",fedm);
   printf(" Number of parameters:         %3i     \n",nparv);
   printf(" Error definition (Fmin + Delta):      %8.3lf  \n",errdef);
   if(fstat==3) {
    printf(" Exact covariance matrix.\n");
   }
   else {
    printf(" No/error with covariance matrix.\n");
    printf(" Error code: %i3\n",fstat);
   }
   printf("\n");


   printf("   Parameter     Value       Error    positive    negative    L_BND    U_BND\n");
   for(i=0;i<nparx;i++) 
    {
        minuit->mnpout(i,pname,pvalue,perror,plbound,pubound,pvari);
        if(pvari>0) // only variable parameters
            {                      
                minuit->mnerrs(i,eplus,eminus,perror,gcorr);
                printf("%2i %10s %10.3le %10.3le %+10.3le %10.3le %8.1le %8.1le\n",
                i,pname.Data(),pvalue,perror,eplus,eminus,plbound,pubound);
            }
   }

   minuit->mnemat(&emat[0][0],MAXPAR);

   for(i=0; i<nparv; i++) {
    for(j=0; j<nparv; j++){ 
     kmat[i][j] = sqrt(emat[i][i]*emat[j][j]); 
      if(kmat[i][j]>1E-80) {
       kmat[i][j] = emat[i][j]/kmat[i][j];
      }
      else kmat[i][j] = 0.0;
     }
   }

   printf("\n");
   printf("Covariance matrix: \n");
   for(i=0; i<nparv; i++) {
    for(j=0; j<nparv; j++){ printf(" %10.3le",emat[i][j]);}
    printf("\n");
   }
   printf("\n");

   printf("Correlation matrix: \n");
   for(i=0; i<nparv; i++) {
     for(j=0; j<nparv; j++){ printf(" %6.3lf",kmat[i][j]);}
     printf("\n");
   }
   printf("\n");
}

void readin(const char* name)
{
   int events;
   double p_plus[3], p_minus[3];

   FILE* fin = fopen(name,"r");        
   if(! fin) {                          
       printf("File can't be opened.\n"); 
       return;                            
   } 

   events = 0;                             
   while(!feof(fin)) {                   
      fscanf(fin,"%lf %lf %lf %lf %lf %lf", 
             &p_plus[0],&p_plus[1],&p_plus[2],&p_minus[0],&p_minus[1],&p_minus[2]);
      costh[events] = p_minus[2]
	  / sqrt(p_minus[0]*p_minus[0]+p_minus[1]*p_minus[1]+p_minus[2]*p_minus[2]);
      events++;                           
      if(events > MAXEVENTS) {           
         printf("Too many events %i\n", MAXEVENTS); 
	 break;                           
      } 
   } 
   events--;                              
   printf("%i Events read. \n", events); 
   fclose(fin);                          

   nevents = events;
}

void fcn(int &npar, double *gin, double &f, double *par, int iflag)
{
   int i;

   double lnL = 0.0;
   for(i=0;i<nevents; i++) {
     lnL += log(3.0/8.0*(1.0+costh[i]*costh[i]) + par[0]*costh[i]);
   }
   f = -lnL;

}

main() 
{ 

   int i;
   int nvor;             
   double rvor;          
   double A;              
   double sigA;         
   double arglist[10];    
   int ierflg = 0;     

   readin("L3.dat");            

   TFile file("asymmetrie.root", "recreate");  
   TH1F hist("hist","costh",20,-1.0,1.0); 
   for(i=0; i<nevents; i++)
      hist.Fill(costh[i]);
   hist.Write();
   file.Close();

// calculate asymmetry by counting
   nvor = 0;
   for(i=0; i<nevents; i++) {
      if(costh[i]>0)
         nvor++;
   }
   rvor = (double)nvor/(double)nevents;
   A = 2.0*rvor - 1.0;   
   sigA = 2.0*sqrt(rvor*(1.0-rvor)/(double)nevents);
   printf("Asymmetry by counting:: \n");
   printf("A = %5.3f +- %5.3f \n", A, sigA);

   TMinuit minuit(1);  

   minuit.SetFCN(fcn);

   arglist[0] = 0.5;
   minuit.mnexcm("SET ERR",arglist,1,ierflg);

   minuit.mnparm(0,"A",0.0,0.1,0,0,ierflg);

   arglist[0] = 500;
   arglist[1] = 1.0;
   minuit.mnexcm("MIGRAD",arglist,2,ierflg);

   arglist[0] = 500;
   arglist[1] = 1.0;
   minuit.mnexcm("MINOS",arglist,2,ierflg);

  results(&minuit);

}




