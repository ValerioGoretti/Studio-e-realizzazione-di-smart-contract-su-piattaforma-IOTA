import 'https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20.sol';

contract TokenViaggio is ERC20{
    
    mapping(address=>bool) private trusted;
    address owner;
    
    constructor()payable ERC20("TokenViaggio", "TV") public {_mint(address(this),500);owner=msg.sender;}
    
    modifier isTheOwner{
        require(msg.sender==owner,"REFUSED");
        _;
    }
    modifier isAtrustedMachine(address addr){
        require(trusted[addr],"The address is not a trusted machine");
        _;
    }
    
    function getBankBalance() public view returns(uint256)
    {
        return balanceOf(address(this));
    }
    
     function setTrustedMachine(address addr) public isTheOwner{trusted[addr]=true;_mint(addr,500);}
    
    function setTrustedAllowance(address owner,uint256 amount) public isAtrustedMachine(msg.sender)
    {
        _approve(owner,msg.sender,amount);
        
    }
     
}
