pragma solidity >=0.4.21 <0.8.0;
contract Ricarica
{
    uint cambioEuroGwei;
    int result;
    address owner;
    constructor() payable{
        cambioEuroGwei=2000000;
        result=0;
        owner=msg.sender;
    }
    modifier isTheOwner(){
        require(msg.sender==owner,"Non sei il proprietario");
        _;
    }
    function eseguiRicarica50() public payable
    {
      msg.sender.send(50*cambioEuroGwei * 1 gwei);
      
    }
     function eseguiRicarica10() public payable
    {
      msg.sender.send(10*cambioEuroGwei * 1 gwei);
      
    }
     function eseguiRicarica20() public payable
    {
      msg.sender.send(20*cambioEuroGwei * 1 gwei);
      
    }
    
    function setCambio(uint value) public isTheOwner{
        cambioEuroGwei=value;
        
    }
    
    function getBalance()public view  returns(uint256){
        
        return address(this).balance;
    }

}