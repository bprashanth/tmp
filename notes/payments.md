India stack:
Aadhaar:
authenticate without putting yourself in front of an officer.
70% of india in villages, contribute 1/3 to gdp, leads to disparity and poverty
only 20% have bank accounts
litrary rates are improving (18% vs 74% since independence)
mobile phones enable ubiquitous connectivity
40billion dollars on schemes for the poor (ration, health etc)
  identity proof for these services
  attached to the village, can't buy 2re rice in tn with bihar ration
  duplicate identities and ghosts - resource capture, 1 person has 4
  common national identity platform for every resident
  can't do this using names and addresses, leather bound books and pens
  uid authority of india
Aadhar
  CID repository
  2fa number and fingerprint
  online auth system returns yes/no
encrypted packet:
  Name, address, gender, dob, photo, 10 fingerprints
compare packet against entire database (compute bound problem)
give her a random number
fingerprint biometric:
  minutae points
  x,y data coordinates form a template
  ridges go bald in rural areas, fingerprints don't work well enough
  iris has a digital template and false accept rates are 0
  5mb packets to CIDR
22 languages with seperate scripts
federal system lends to state government to implement projects
state gov is a registrar
private companies buy enrollment kits
training agencies
devices to record bio information and certification agencies
bank is not in the village
opportunity cost of a wasted day, transportation cost, 20% of income burnt
people need to withdraw cash in their own villages
microatms
  smartphone + fingerprint scanner
compute:
  million registrations per day, 300 trillion bio metric batches
  in a DC < 5000 sq feet
  storage 5pb - 20 pb
egovernance projects in india


open source deployments:
cidr:
  hadoop distributed fs to store encrypted packets
  mysql, doesn't scale horizontally
  mongodb with solr indexing to do text searches
  several stages of structure validation
  rabbit inbetween stages
  linux and java
  spring aop framework
  commodity blade servers
  zookeeper for orchestration
biometric system:
  3 provider
  ref no and a bunch of templates
privacy and security:
  number itself is random
  data encrypted at source (pki rsa)
  biometrically signed by operator
    with his biometric data - nonrepudiation and penal provisions under law
  data partitoined across multiple security zones
  uid is opaque system - output is y/n
  enrollment system is batch, auth system is live online
grc provider
  governance risk and compliance
  end to end security
national rural employment guarantee scheme:
  guarantee 100 days of work to poor families at minimum rate


Trends:
1. Transactions:
bank: low volume, high value, high cost
ph: high volume, high value, low cost

electronic vs paper clearing
IMPS - npci non profit section 25 setup and owned by banks (visa and mastercard).
instant credit, remmittance
imps overtakes debit
march 2017 overtakes both

2. credentials
proprietary to open
2fa - financial provider responsible for card
phone replaces card
  use phone with pin
  use phone with aadhaar auth

3. switching costs are coming down
dual sim - one sim for incoming calls, deals; one sim for voice one for data.
interoperability of payments
savings rate in india are de regulated (4%, dbs bank from sing 7%)
lending rates

4. lending goes form uniform rates to individual pricing of risk
digital footprints
lending against assest to lending against cash flow
in a services economy they don't have assets they have cash flow

5. business models
fees vs data
no license, digital advertising market

6. psu banks have a shrinking market share
telecom - bsnl and mtnl have lost 70% market to mobiles
airlines - national carrier lots 25% in 4 years
psu banks are losing market share according to rbi, down to 63% by 2025
indian banking sector market cap is coming down - sbi 170k crore, all psu banks 180k crore
private banks - axis, kotak etc 700k cr

