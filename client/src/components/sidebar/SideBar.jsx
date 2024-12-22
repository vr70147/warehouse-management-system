import { useState } from 'react';
import { NavLink } from 'react-router-dom';
import inventory from '../../assets/inventory.svg';
import dashboard from '../../assets/dashboard.svg';
import shipment from '../../assets/shipment.svg';
import order from '../../assets/order.svg';
import analytics from '../../assets/analytics.svg';
import support from '../../assets/support.svg';
import setting from '../../assets/setting.svg';
import { Button } from '../ui/button';
import {
  ChevronLeft,
  ChevronRight,
  ChevronDown,
  ChevronUp,
} from 'lucide-react';

export default function SideBar() {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isAnalyticsExpanded, setIsAnalyticsExpanded] = useState(false);

  const sections = [
    { name: 'Dashboard', icon: dashboard, path: '/dashboard' },
    { name: 'Inventory', icon: inventory, path: '/inventory' },
    { name: 'Shipments', icon: shipment, path: '/shipment' },
    { name: 'Orders', icon: order, path: '/order' },
    {
      name: 'Analytics',
      icon: analytics,
      children: [
        { name: 'Sales Overview', path: '/analytics/sales-overview' },
        { name: 'Customer Insights', path: '/analytics/customer-insights' },
        { name: 'Product Performance', path: '/analytics/product-performance' },
        { name: 'Shipment Analytics', path: '/analytics/shipment-analytics' },
        { name: 'Financial Reports', path: '/analytics/financial-reports' },
      ],
    },
    { name: 'Support', icon: support, path: '/support' },
    { name: 'Settings', icon: setting, path: '/setting' },
  ];

  return (
    <div
      className={`${
        isCollapsed ? 'w-20' : 'w-52'
      } bg-white dark:bg-gray-800 h-screen shadow-lg transition-all duration-500 overflow-hidden`}
    >
      {/* Toggle Button */}
      <div className="flex justify-end p-2">
        <Button
          variant="outline"
          size="icon"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          {isCollapsed ? <ChevronRight /> : <ChevronLeft />}
        </Button>
      </div>
      {/* Navigation Links */}
      <nav className="flex flex-col space-y-2 p-2">
        {sections.map((section) =>
          section.children ? (
            <div key={section.name} className="flex flex-col">
              {/* Parent Link */}
              <button
                onClick={() => setIsAnalyticsExpanded(!isAnalyticsExpanded)}
                className="group flex items-center gap-3 p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-900 transition-all"
              >
                <div className="flex items-center justify-center h-10 w-10 flex-shrink-0">
                  <img
                    src={section.icon}
                    alt={`${section.name} icon`}
                    className="h-6 w-6 dark:invert group-hover:scale-125 transition-transform duration-300"
                  />
                </div>
                <span
                  className={`text-sm font-medium transition-all duration-300 ease-in-out group-hover:scale-125 ${
                    isCollapsed ? 'opacity-0 scale-90' : 'opacity-100 scale-100'
                  }`}
                >
                  {section.name}
                </span>
                <span className="ml-auto">
                  {isAnalyticsExpanded ? <ChevronUp /> : <ChevronDown />}
                </span>
              </button>
              {/* Child Links */}
              <div
                className={`overflow-hidden transition-all duration-300 ${
                  isAnalyticsExpanded
                    ? 'max-h-screen opacity-100'
                    : 'max-h-0 opacity-0'
                }`}
              >
                {section.children.map((child) => (
                  <NavLink
                    key={child.name}
                    to={child.path}
                    className={({ isActive }) =>
                      `group flex items-center gap-3 p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-900 transition-all ${
                        isActive ? 'dark:bg-gray-900 bg-gray-200' : ''
                      }`
                    }
                  >
                    <span
                      className={`text-sm font-medium transition-all duration-300 ease-in-out ${
                        isCollapsed
                          ? 'opacity-0 scale-90'
                          : 'opacity-100 scale-100'
                      }`}
                    >
                      {child.name}
                    </span>
                  </NavLink>
                ))}
              </div>
            </div>
          ) : (
            <NavLink
              key={section.name}
              to={section.path}
              className={({ isActive }) =>
                `group flex items-center gap-3 p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-900 hover: transition-all ${
                  isActive ? 'dark:bg-gray-900 bg-gray-200' : ''
                }`
              }
            >
              <div className="flex items-center justify-center h-10 w-10 flex-shrink-0">
                <img
                  src={section.icon}
                  alt={`${section.name} icon`}
                  className="h-6 w-6 dark:invert group-hover:scale-125 transition-transform duration-300"
                />
              </div>
              <span
                className={`text-sm font-medium transition-all duration-300 ease-in-out group-hover:scale-125 ${
                  isCollapsed ? 'opacity-0 scale-90' : 'opacity-100 scale-100'
                }`}
              >
                {section.name}
              </span>
            </NavLink>
          )
        )}
      </nav>
    </div>
  );
}
